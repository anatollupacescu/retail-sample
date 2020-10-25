package stock

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
)

type (
	entity struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Qty  int    `json:"qty"`
	}
	singleResponse struct {
		Data entity `json:"data"`
	}
	collectionResponse struct {
		Data []entity `json:"data"`
	}
)

func internalServerErrorMsg() string {
	return http.StatusText(http.StatusInternalServerError)
}

func GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	uc, err := useCase(r)

	if err != nil {
		http.Error(w, internalServerErrorMsg(), http.StatusBadRequest)
		return
	}

	entries, err := uc.CurrentStock()

	var response collectionResponse
	response.Data = make([]entity, 0, len(entries))

	for i := range entries {
		entry := entries[i]

		response.Data = append(response.Data, entity{
			ID:   entry.ID,
			Name: entry.Name,
			Qty:  entry.Qty,
		})
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, internalServerErrorMsg(), http.StatusInternalServerError)
	}
}

func Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	uc, err := useCase(r)

	if err != nil {
		http.Error(w, internalServerErrorMsg(), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	itemID := vars["itemID"]

	pos, err := uc.Position(itemID)

	switch err {
	case nil:
		break
	case inventory.ErrItemNotFound:
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	default:
		http.Error(w, internalServerErrorMsg(), http.StatusInternalServerError)
		return
	}

	var response = singleResponse{
		Data: entity{
			ID:   pos.ID,
			Qty:  pos.Qty,
			Name: pos.Name,
		},
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, internalServerErrorMsg(), http.StatusInternalServerError)
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
		http.Error(w, internalServerErrorMsg(), http.StatusBadRequest)
		return
	}

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var body provisionPayload

	if err = d.Decode(&body); err != nil {
		http.Error(w, internalServerErrorMsg(), http.StatusInternalServerError)

		return
	}

	vars := mux.Vars(r)
	itemID := vars["itemID"]

	pos, err := uc.Provision(itemID, body.Qty)

	switch err {
	case nil:
		break
	case inventory.ErrItemNotFound:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		http.Error(w, internalServerErrorMsg(), http.StatusInternalServerError)
		return
	}

	var response = provisionResponse{
		Data: provisionResponseData{
			ID:  pos.ID,
			Qty: pos.Qty,
		},
	}

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, internalServerErrorMsg(), http.StatusInternalServerError)
	}
}

func GetProvisionLog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	uc, err := useCase(r)

	if err != nil {
		http.Error(w, internalServerErrorMsg(), http.StatusBadRequest)
		return
	}

	pl, err := uc.ProvisionLog()

	if err != nil {
		http.Error(w, internalServerErrorMsg(), http.StatusBadRequest)
	}

	type provisionLogEntity struct {
		Time time.Time `json:"time"`
		ID   int       `json:"id"`
		Qty  int       `json:"qty"`
	}

	var response struct {
		Data []provisionLogEntity `json:"data"`
	}

	response.Data = make([]provisionLogEntity, 0)

	for i := range pl {
		in := pl[i]

		response.Data = append(response.Data, provisionLogEntity{
			ID:  in.ID,
			Qty: in.Qty,
		})
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, internalServerErrorMsg(), http.StatusInternalServerError)
	}
}
