package recipe

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/hlog"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/usecase"
	domain "github.com/anatollupacescu/retail-sample/domain/retail/recipe"
)

type updatePayload struct {
	Enabled bool `json:"enabled"`
}

var ErrParseItemID = errors.New("could not parse item ID")

func newUpdateStatusDTO(r *http.Request) (usecase.UpdateStatusDTO, error) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var requestBody updatePayload

	if err := d.Decode(&requestBody); err != nil {
		hlog.FromRequest(r).Err(err).Msg("parse 'update recipe status' payload")
		return usecase.UpdateStatusDTO{}, err
	}

	vars := mux.Vars(r)
	recipeID := vars["recipeID"]

	id, err := strconv.Atoi(recipeID)

	if err != nil {
		hlog.FromRequest(r).Err(err).Msg("parse 'update status' itemID")
		return usecase.UpdateStatusDTO{}, ErrParseItemID
	}

	dto := usecase.UpdateStatusDTO{
		RecipeID: id,
		Enabled:  requestBody.Enabled,
	}

	return dto, nil
}

type createPayload struct {
	Name  string      `json:"name"` // pointer so we can test for field absence
	Items map[int]int `json:"items"`
}

func newCreateDTO(r *http.Request) (usecase.CreateRecipeDTO, error) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var requestBody createPayload

	if err := d.Decode(&requestBody); err != nil {
		hlog.FromRequest(r).Err(err).Msg("parse 'create' payload")
		return usecase.CreateRecipeDTO{}, err
	}

	var ingredients = make([]domain.Ingredient, 0, len(requestBody.Items))

	for id, qty := range requestBody.Items {
		ingredients = append(ingredients, domain.Ingredient{
			ID:  id,
			Qty: qty,
		})
	}

	dto := usecase.CreateRecipeDTO{
		Name:        domain.Name(requestBody.Name),
		Ingredients: ingredients,
	}

	return dto, nil
}
