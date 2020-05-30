package recipe

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/middleware"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/recipe"
)

type (
	webApp struct {
		wrapper wrapper
		logger  middleware.Logger
	}

	item struct {
		ID  int `json:"id"`
		Qty int `json:"qty"`
	}

	entity struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		Items   []item `json:"items"`
		Enabled bool   `json:"enabled"`
	}
)

var internalServerError = "internal server error"

func (a *webApp) create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var requestBody struct {
		Name  string      `json:"name"` // pointer so we can test for field absence
		Items map[int]int `json:"items"`
	}

	if err := d.Decode(&requestBody); err != nil {
		a.logger.Log("action", "decode request payload", "error", err, "method", "recipe.create")
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

	re, err := a.wrapper.create(recipeName, ingredients)

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
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	var response = struct {
		Data entity `json:"data"`
	}{
		Data: entity{
			ID:      int(re.ID),
			Name:    string(re.Name),
			Items:   toItems(re.Ingredients),
			Enabled: re.Enabled,
		},
	}

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		a.logger.Log("action", "encode response", "error", err, "method", "recipe.create")
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}
}

func (a *webApp) getAll(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	all, err := a.wrapper.getAll()

	if err != nil {
		http.Error(w, internalServerError, http.StatusBadRequest)
		return
	}

	var response struct {
		Data []entity `json:"data"`
	}

	response.Data = make([]entity, 0)

	for _, r := range all {
		response.Data = append(response.Data, entity{
			ID:      int(r.ID),
			Name:    string(r.Name),
			Items:   toItems(r.Ingredients),
			Enabled: r.Enabled,
		})
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		a.logger.Log("action", "encode response", "error", err, "method", "recipe.getAll")
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}
}

func toItems(i []recipe.Ingredient) (items []item) {
	for _, ri := range i {
		items = append(items, item{
			ID:  int(ri.ID),
			Qty: int(ri.Qty),
		})
	}

	return
}

func (a *webApp) get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	rid := vars["recipeID"]

	id, err := strconv.Atoi(rid)

	if err != nil {
		a.logger.Log("action", "parse id", "error", err, "method", "recipe.get")
		http.Error(w, "could not parse id", http.StatusBadRequest)
		return
	}

	recipeID := recipe.ID(id)

	rcp, err := a.wrapper.get(recipeID)

	switch err {
	case nil:
		break
	case recipe.ErrRecipeNotFound:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	var response = struct {
		Data entity `json:"data"`
	}{
		Data: entity{
			ID:      id,
			Name:    string(rcp.Name),
			Items:   toItems(rcp.Ingredients),
			Enabled: rcp.Enabled,
		},
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		a.logger.Log("action", "encode response", "error", err, "method", "recipe.get")
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}
}

func (a *webApp) update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	rid := vars["recipeID"]

	id, err := strconv.Atoi(rid)

	if err != nil {
		a.logger.Log("action", "parse id", "error", err, "method", "recipe.update")
		http.Error(w, "could not parse id", http.StatusBadRequest)
		return
	}

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var requestBody struct {
		Enabled bool `json:"enabled"`
	}

	if err := d.Decode(&requestBody); err != nil {
		a.logger.Log("action", "decode request payload", "error", err, "method", "recipe.update")
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	re, err := a.wrapper.setStatus(id, requestBody.Enabled)

	switch err {
	case nil:
		break
	case recipe.ErrRecipeNotFound:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	var response = struct {
		Data entity `json:"data"`
	}{
		Data: entity{
			ID:      int(re.ID),
			Name:    string(re.Name),
			Items:   toItems(re.Ingredients),
			Enabled: re.Enabled,
		},
	}

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		a.logger.Log("action", "encode response", "error", err, "method", "recipe.update")
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}
}