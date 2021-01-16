package inventory

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	usecase "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/machine"
	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/machine/inventory"
)

type updatePayload struct {
	Enabled bool `json:"enabled"`
}

func Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var requestPayload updatePayload
	if err := d.Decode(&requestPayload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uc, err := inventory.New(r.Context())

	if err != nil {
		httpServerError(w)
		return
	}

	vars := mux.Vars(r)
	id := vars["itemID"]

	item, err := uc.UpdateStatus(id, requestPayload.Enabled)

	switch {
	case err == nil:
	case errors.Is(err, usecase.ErrNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	case errors.Is(err, usecase.ErrBadRequest):
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		httpServerError(w)
		return
	}

	w.WriteHeader(http.StatusAccepted)

	response := toSingleResponse(item)
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		httpServerError(w)
	}
}

func Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var body createPayload

	if err := d.Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uc, err := inventory.New(r.Context())

	if err != nil {
		httpServerError(w)
		return
	}

	newItem, err := uc.Create(body.Name)

	switch {
	case err == nil:
	case errors.Is(err, usecase.ErrBadRequest):
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		httpServerError(w)
		return
	}

	w.WriteHeader(http.StatusCreated)

	response := toSingleResponse(newItem)
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		httpServerError(w)
	}
}

type createPayload struct {
	Name string `json:"name"`
}

func Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	uc, err := inventory.New(r.Context())

	if err != nil {
		httpServerError(w)
		return
	}

	vars := mux.Vars(r)
	id := vars["itemID"]

	item, err := uc.GetByID(id)

	switch {
	case err == nil:
	case errors.Is(err, usecase.ErrBadRequest):
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	case errors.Is(err, usecase.ErrNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	default:
		httpServerError(w)
		return
	}

	response := toSingleResponse(item)
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		httpServerError(w)
	}
}

func GetAll(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			httpServerError(w)
		}
	}()

	w.Header().Set("Content-Type", "application/json")

	uc, err := inventory.New(r.Context())

	if err != nil {
		return
	}

	items, err := uc.GetAll()

	if err != nil {
		return
	}

	response := toCollectionResponse(items)
	err = json.NewEncoder(w).Encode(response)
}
