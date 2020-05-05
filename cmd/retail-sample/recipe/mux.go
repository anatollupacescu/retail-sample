package recipe

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/types"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/recipe"
)

type webApp struct {
	wrapper wrapper
	logger  types.Logger
}

var internalServerError = "internal server error"

func (a *webApp) CreateRecipe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var requestBody struct {
		Name  string      `json:"name"` // pointer so we can test for field absence
		Items map[int]int `json:"items"`
	}

	if err := d.Decode(&requestBody); err != nil {
		a.logger.Log("action", "decode request payload", "error", err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	var ingredients []recipe.Ingredient

	for id, qty := range requestBody.Items {
		ingredients = append(ingredients, recipe.Ingredient{
			ID:  id,
			Qty: qty,
		})
	}

	var recipeName = recipe.Name(requestBody.Name)

	recipeID, err := a.wrapper.Add(recipeName, ingredients)

	switch err {
	case nil:
		break
	case recipe.ErrEmptyName:
		fallthrough
	case recipe.ErrIgredientNotFound:
		fallthrough
	case recipe.ErrNoIngredients:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		a.logger.Log("action", "call application", "error", err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	type descriptorEntity map[recipe.Name]recipe.ID

	var response = struct {
		Data descriptorEntity `json:"data"`
	}{
		Data: descriptorEntity{
			recipeName: recipeID,
		},
	}

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}
}

func (a *webApp) ListRecipes(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	list, err := a.wrapper.List()

	if err != nil {
		a.logger.Log("action", "call application", "error", err)
		http.Error(w, internalServerError, http.StatusBadRequest)
	}

	type recipe struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Items []item `json:"items"`
	}

	var response struct {
		Data []recipe `json:"data"`
	}

	response.Data = make([]recipe, 0)

	for _, r := range list {
		response.Data = append(response.Data, recipe{
			ID:    int(r.ID),
			Name:  string(r.Name),
			Items: toItems(r.Ingredients),
		})
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		a.logger.Log("action", "encode response", "error", err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}
}

type item struct {
	Id  int `json:"id"`
	Qty int `json:"qty"`
}

func toItems(i []recipe.Ingredient) (items []item) {
	for _, ri := range i {
		items = append(items, item{
			Id:  int(ri.ID),
			Qty: int(ri.Qty),
		})
	}

	return
}

func (a *webApp) GetRecipe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	rid := vars["recipeID"]

	i, _ := strconv.Atoi(rid)

	recipeID := recipe.ID(i)

	rcp, err := a.wrapper.Get(recipeID)

	switch err {
	case nil:
		break
	case recipe.ErrRecipeNotFound:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		a.logger.Log("action", "call application", "error", err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	var response = struct {
		Name  string `json:"name"`
		Items []item `json:"items"`
	}{
		Name:  string(rcp.Name),
		Items: toItems(rcp.Ingredients),
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		a.logger.Log("action", "encode response", "error", err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}
}
