package recipe

import (
	"encoding/json"
	"net/http"

	"github.com/anatollupacescu/retail-sample/domain/retail/recipe"
)

var serverErr = http.StatusText(http.StatusInternalServerError)

func Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	uc, err := useCase(r)

	if err != nil {
		http.Error(w, serverErr, http.StatusInternalServerError)
		return
	}

	dto, err := toCreateDTO(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	re, err := uc.Create(dto)

	switch err {
	case nil:
		break
	case
		recipe.ErrEmptyName,
		recipe.ErrIgredientNotFound,
		recipe.ErrQuantityNotProvided,
		recipe.ErrNoIngredients:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		http.Error(w, serverErr, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	response := toResponse(re)
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, serverErr, http.StatusInternalServerError)
	}
}

func GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	all, err := getAll(r)

	if err != nil {
		http.Error(w, serverErr, http.StatusInternalServerError)
		return
	}

	response := toCollectionResponse(all)
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, serverErr, http.StatusInternalServerError)
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
		http.Error(w, serverErr, http.StatusInternalServerError)
		return
	}

	response := toResponse(rcp)
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, serverErr, http.StatusInternalServerError)
	}
}

func Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	uc, err := useCase(r)

	if err != nil {
		http.Error(w, serverErr, http.StatusInternalServerError)
		return
	}

	dto, err := toUpdateStatusDTO(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedRecipe, err := uc.Update(dto)

	switch err {
	case nil:
		break
	case recipe.ErrRecipeNotFound:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		http.Error(w, serverErr, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	response := toResponse(updatedRecipe)
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, serverErr, http.StatusInternalServerError)
	}
}
