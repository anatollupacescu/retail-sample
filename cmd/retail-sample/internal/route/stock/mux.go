package stock

import (
	"encoding/json"
	"net/http"

	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
)

var serverErr = http.StatusText(http.StatusInternalServerError)

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

	position, err := getByID(r)

	switch err {
	case nil:
		break
	case inventory.ErrItemNotFound, ErrBadItemID:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		http.Error(w, serverErr, http.StatusInternalServerError)
		return
	}

	response := toResponse(position)
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, serverErr, http.StatusInternalServerError)
	}
}

type provisionPayload struct {
	Qty int `json:"qty"`
}

type provisionResponseData struct {
	ID  int `json:"id"`
	Qty int `json:"qty"`
}

type provisionResponse struct {
	Data provisionResponseData `json:"data"`
}

func Provision(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	uc, err := useCase(r)

	if err != nil {
		http.Error(w, serverErr, http.StatusInternalServerError)
		return
	}

	dto, err := toProvisionDTO(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pos, err := uc.Provision(dto)

	switch err {
	case nil:
		break
	case inventory.ErrItemNotFound:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		http.Error(w, serverErr, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	response := toResponse(pos)
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, serverErr, http.StatusInternalServerError)
	}
}

func GetProvisionLog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	pl, err := getProvisionLog(r)

	if err != nil {
		http.Error(w, serverErr, http.StatusInternalServerError)
	}

	response := toProvisionLog(pl)
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, serverErr, http.StatusInternalServerError)
	}
}
