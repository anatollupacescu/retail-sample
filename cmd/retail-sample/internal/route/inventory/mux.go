package inventory

import (
	"encoding/json"
	"net/http"

	usecase "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/usecase/inventory"

	"github.com/gorilla/mux"

	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
)

type (
	entity struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		Enabled bool   `json:"enabled"`
	}
	singleResponse struct {
		Data entity `json:"data"`
	}
	collectionResponse struct {
		Data []entity `json:"data"`
	}
)

type updatePayload struct {
	Enabled bool `json:"enabled"`
}

func Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	uc, err := useCase(r)

	if err != nil {
		http.Error(w, internalServerErrorMsg(), http.StatusInternalServerError)

		return
	}

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var requestPayload updatePayload

	if err := d.Decode(&requestPayload); err != nil {
		http.Error(w, "parse body", http.StatusBadRequest)

		return
	}

	vars := mux.Vars(r)
	rid := vars["itemID"]

	item, err := uc.UpdateStatus(rid, requestPayload.Enabled)

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
			ID:      item.ID,
			Name:    item.Name,
			Enabled: item.Enabled,
		},
	}

	w.WriteHeader(http.StatusAccepted)

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, internalServerErrorMsg(), http.StatusInternalServerError)
	}
}

type createPayload struct {
	Name string `json:"name"`
}

func Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	uc, err := useCase(r)

	if err != nil {
		http.Error(w, internalServerErrorMsg(), http.StatusInternalServerError)
		return
	}

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var body createPayload

	if err := d.Decode(&body); err != nil {
		http.Error(w, "could not parse body", http.StatusBadRequest)
		return
	}

	newItem, err := uc.Create(body.Name)

	switch err {
	case nil:
		break
	case inventory.ErrEmptyName, inventory.ErrDuplicateName:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		http.Error(w, internalServerErrorMsg(), http.StatusInternalServerError)
		return
	}

	var response = singleResponse{
		Data: entity{
			ID:      newItem.ID,
			Name:    newItem.Name,
			Enabled: newItem.Enabled,
		},
	}

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		// logger.Log("action", "encode response", "error", err, "method", "inventory.create")
		http.Error(w, internalServerErrorMsg(), http.StatusInternalServerError)
	}
}

func Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	uc, err := useCase(r)

	if err != nil {
		http.Error(w, internalServerErrorMsg(), http.StatusInternalServerError)

		return
	}

	vars := mux.Vars(r)
	itemID := vars["itemID"]

	item, err := uc.GetByID(itemID)

	switch err {
	case nil:
		break
	case usecase.ErrBadItemID:
		http.Error(w, "could not parse id: "+itemID, http.StatusBadRequest)
		return
	case inventory.ErrItemNotFound:
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	default:
		http.Error(w, internalServerErrorMsg(), http.StatusInternalServerError)
		return
	}

	var response = singleResponse{
		Data: entity{
			ID:      item.ID,
			Name:    item.Name,
			Enabled: item.Enabled,
		},
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, internalServerErrorMsg(), http.StatusInternalServerError)
	}
}

func GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	uc, err := useCase(r)

	if err != nil {
		http.Error(w, internalServerErrorMsg(), http.StatusInternalServerError)

		return
	}

	all, err := uc.GetAll()

	switch err {
	case nil:
		break
	default:
		http.Error(w, internalServerErrorMsg(), http.StatusInternalServerError)
		return
	}

	var response = collectionResponse{
		Data: make([]entity, 0, len(all)),
	}

	for _, item := range all {
		response.Data = append(response.Data, entity{
			ID:      item.ID,
			Name:    item.Name,
			Enabled: item.Enabled,
		})
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, internalServerErrorMsg(), http.StatusInternalServerError)
	}
}

func internalServerErrorMsg() string {
	return http.StatusText(http.StatusInternalServerError)
}
