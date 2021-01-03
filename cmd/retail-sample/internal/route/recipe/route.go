package recipe

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	usecase "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/machine/recipe"
	"github.com/anatollupacescu/retail-sample/domain/retail/recipe"
)

type createPayload struct {
	Name  string      `json:"name"`
	Items map[int]int `json:"items"`
}

func Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var requestBody createPayload

	if err := d.Decode(&requestBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var ingredients = make([]recipe.InventoryItem, 0, len(requestBody.Items))

	for id, qty := range requestBody.Items {
		ingredients = append(ingredients, recipe.InventoryItem{
			ID:  id,
			Qty: qty,
		})
	}

	uc, err := usecase.New(r.Context())
	if err != nil {
		httpServerError(w)
		return
	}

	re, err := uc.Create(requestBody.Name, ingredients)

	switch err {
	case nil:
		break
	case
		recipe.ErrEmptyName,
		recipe.ErrIngredientNotFound,
		recipe.ErrQuantityNotProvided,
		recipe.ErrNoIngredients:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		httpServerError(w)
		return
	}

	w.WriteHeader(http.StatusCreated)

	response := toResponse(re)
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		httpServerError(w)
	}
}

func GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	all, err := getAll(r)

	if err != nil {
		httpServerError(w)
		return
	}

	response := toCollectionResponse(all)
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		httpServerError(w)
	}
}

func Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rcp, err := getByID(r)

	switch err {
	case nil:
		break
	case recipe.ErrRecipeNotFound, ErrBadItemID:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		httpServerError(w)
		return
	}

	response := toResponse(rcp)
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		httpServerError(w)
	}
}

type updatePayload struct {
	Enabled bool `json:"enabled"`
}

func Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var requestBody updatePayload

	if err := d.Decode(&requestBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	recipeID := vars["recipeID"]

	id, err := strconv.Atoi(recipeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uc, err := usecase.New(r.Context())
	if err != nil {
		httpServerError(w)
		return
	}

	updatedRecipe, err := uc.UpdateStatus(id, requestBody.Enabled)

	switch err {
	case nil:
		break
	case recipe.ErrRecipeNotFound:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		httpServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)

	response := toResponse(updatedRecipe)
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		httpServerError(w)
	}
}
