package recipe

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/usecase/recipe"
	domain "github.com/anatollupacescu/retail-sample/domain/retail/recipe"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/hlog"
)

type updatePayload struct {
	Enabled bool `json:"enabled"`
}

var ErrParseItemID = errors.New("could not parse item ID")

func toUpdateStatusDTO(r *http.Request) (recipe.UpdateStatusDTO, error) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var requestBody updatePayload

	if err := d.Decode(&requestBody); err != nil {
		hlog.FromRequest(r).Err(err).Msg("parse 'update recipe status' payload")
		return recipe.UpdateStatusDTO{}, err
	}

	vars := mux.Vars(r)
	recipeID := vars["recipeID"]

	id, err := strconv.Atoi(recipeID)

	if err != nil {
		hlog.FromRequest(r).Err(err).Msg("parse 'update status' itemID")
		return recipe.UpdateStatusDTO{}, ErrParseItemID
	}

	dto := recipe.UpdateStatusDTO{
		RecipeID: id,
		Enabled:  requestBody.Enabled,
	}

	return dto, nil
}

type createPayload struct {
	Name  string      `json:"name"` // pointer so we can test for field absence
	Items map[int]int `json:"items"`
}

func toCreateDTO(r *http.Request) (recipe.CreateDTO, error) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var requestBody createPayload

	if err := d.Decode(&requestBody); err != nil {
		hlog.FromRequest(r).Err(err).Msg("parse 'create' payload")
		return recipe.CreateDTO{}, err
	}

	var ingredients = make([]domain.Ingredient, 0, len(requestBody.Items))

	for id, qty := range requestBody.Items {
		ingredients = append(ingredients, domain.Ingredient{
			ID:  id,
			Qty: qty,
		})
	}

	dto := recipe.CreateDTO{
		Name:        domain.Name(requestBody.Name),
		Ingredients: ingredients,
	}

	return dto, nil
}

// func toIngredientCollection(items []item) []domain.Ingredient {
// 	ingredients := make([]domain.Ingredient, 0, len(items))

// 	for i := range items {
// 		currentItem := items[i]

// 		ingredients = append(ingredients, domain.Ingredient{
// 			ID:  currentItem.ID,
// 			Qty: currentItem.Qty,
// 		})
// 	}

// 	return ingredients
// }
