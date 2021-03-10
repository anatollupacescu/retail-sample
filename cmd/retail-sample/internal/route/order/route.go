package order

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	usecase "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/machine"
	order "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/machine/order"
)

type createPayload struct {
	RecipeID int `json:"id"`
	Count    int `json:"qty"`
}

func Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var payload createPayload

	if err := d.Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uc, err := order.New(r.Context())

	if err != nil {
		httpServerError(w)
		return
	}

	newOrder, err := uc.Create(payload.RecipeID, payload.Count)

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

	w.WriteHeader(http.StatusCreated)

	var response = toSingleResponse(newOrder)
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		httpServerError(w)
	}
}

func Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	uc, err := order.New(r.Context())

	if err != nil {
		httpServerError(w)
		return
	}

	vars := mux.Vars(r)
	orderID := vars["orderID"]

	ordr, err := uc.GetByID(orderID)

	switch err {
	case nil:
	case usecase.ErrNotFound:
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	case usecase.ErrBadRequest:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		httpServerError(w)
		return
	}

	response := toSingleResponse(ordr)
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

	uc, err := order.New(r.Context())

	if err != nil {
		return
	}

	orders, err := uc.GetAll()

	if err != nil {
		return
	}

	response := toCollectionResponse(orders)
	err = json.NewEncoder(w).Encode(response)
}
