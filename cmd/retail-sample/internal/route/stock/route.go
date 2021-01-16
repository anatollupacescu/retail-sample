package stock

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	usecase "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/machine"
	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/machine/stock"
)

func GetAll(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			httpServerError(w)
		}
	}()

	w.Header().Set("Content-Type", "application/json")

	uc, err := stock.New(r.Context())

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

	uc, err := stock.New(r.Context())

	if err != nil {
		httpServerError(w)
		return
	}

	vars := mux.Vars(r)
	itemID := vars["itemID"]

	position, err := uc.GetByID(itemID)

	switch {
	case err == nil:
	case errors.Is(err, usecase.ErrNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	default:
		httpServerError(w)
		return
	}

	response := toResponse(position)
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		httpServerError(w)
	}
}

type provisionPayload struct {
	Qty int `json:"qty"`
}

func Provision(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var body provisionPayload

	if err := d.Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uc, err := stock.New(r.Context())

	if err != nil {
		httpServerError(w)
		return
	}

	vars := mux.Vars(r)
	itemID := vars["itemID"]

	pos, err := uc.Provision(itemID, body.Qty)

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

	w.WriteHeader(http.StatusCreated)

	response := toResponse(pos)
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		httpServerError(w)
	}
}

func GetProvisionLog(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			httpServerError(w)
		}
	}()

	w.Header().Set("Content-Type", "application/json")

	uc, err := stock.New(r.Context())

	if err != nil {
		return
	}

	pl, err := uc.GetProvisionLog()

	if err != nil {
		return
	}

	response := toProvisionLog(pl)

	err = json.NewEncoder(w).Encode(response)
}
