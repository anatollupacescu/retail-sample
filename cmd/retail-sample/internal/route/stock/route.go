package stock

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/machine/stock"
	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
)

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

	position, err := getByID(r)

	switch err {
	case nil:
		break
	case inventory.ErrItemNotFound, ErrBadItemID:
		http.Error(w, err.Error(), http.StatusBadRequest)
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

	vars := mux.Vars(r)
	itemID := vars["itemID"]

	id, err := strconv.Atoi(itemID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uc, err := stock.New(r.Context())

	if err != nil {
		httpServerError(w)
		return
	}

	pos, err := uc.Provision(id, body.Qty)

	switch err {
	case nil:
		break
	case inventory.ErrItemNotFound:
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

var ErrParseItemID = errors.New("could not parse item ID")

func GetProvisionLog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	pl, err := getProvisionLog(r)

	if err != nil {
		httpServerError(w)
	}

	response := toProvisionLog(pl)
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		httpServerError(w)
	}
}
