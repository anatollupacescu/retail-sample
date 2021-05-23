package recipe

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	machine "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/machine"
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

	switch {
	case err == nil:
	case errors.Is(err, machine.ErrNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	case errors.Is(err, machine.ErrBadRequest):
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

	var err error

	defer func() {
		if err != nil {
			httpServerError(w)
		}
	}()

	uc, err := usecase.New(r.Context())

	if err != nil {
		return
	}

	all, err := uc.GetAll()

	if err != nil {
		return
	}

	response := toCollectionResponse(all)
	err = json.NewEncoder(w).Encode(response)
}

func Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	uc, err := usecase.New(r.Context())

	if err != nil {
		httpServerError(w)
		return
	}

	vars := mux.Vars(r)
	recipeID := vars["recipeID"]

	rcp, err := uc.GetByID(recipeID)

	switch {
	case err == nil:
	case errors.Is(err, machine.ErrNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	case errors.Is(err, machine.ErrBadRequest):
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

	uc, err := usecase.New(r.Context())

	if err != nil {
		httpServerError(w)
		return
	}

	updatedRecipe, err := uc.UpdateStatus(recipeID, requestBody.Enabled)

	switch {
	case err == nil:
	case errors.Is(err, machine.ErrNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	case errors.Is(err, machine.ErrBadRequest):
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
