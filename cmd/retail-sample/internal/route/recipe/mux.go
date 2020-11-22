package recipe

import (
	"encoding/json"
	"net/http"

	"github.com/anatollupacescu/retail-sample/domain/retail/recipe"
)

func httpServerError(w http.ResponseWriter) {
	status := http.StatusInternalServerError
	statusText := http.StatusText(status)
	http.Error(w, statusText, http.StatusInternalServerError)
}

func Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	uc, err := useCase(r)

	if err != nil {
		httpServerError(w)
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

func Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	uc, err := useCase(r)

	if err != nil {
		httpServerError(w)
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
