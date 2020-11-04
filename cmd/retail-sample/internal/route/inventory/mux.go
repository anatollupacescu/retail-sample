package inventory

import (
	"encoding/json"
	"net/http"

	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
)

var serverError = http.StatusText(http.StatusInternalServerError)

func Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	uc, err := useCase(r)

	if err != nil {
		http.Error(w, serverError, http.StatusInternalServerError)
		return
	}

	dto, err := toUpdateStatusDTO(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	item, err := uc.UpdateStatus(dto)

	switch err {
	case nil:
		break
	case inventory.ErrItemNotFound:
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	default:
		http.Error(w, serverError, http.StatusInternalServerError)
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
		http.Error(w, serverError, http.StatusInternalServerError)
	}
}

func Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	uc, err := useCase(r)

	if err != nil {
		http.Error(w, serverError, http.StatusInternalServerError)
		return
	}

	dto, err := toCreateDTO(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newItem, err := uc.Create(dto)

	switch err {
	case nil:
		break
	case inventory.ErrEmptyName, inventory.ErrDuplicateName:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		http.Error(w, serverError, http.StatusInternalServerError)
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
		http.Error(w, serverError, http.StatusInternalServerError)
	}
}

func Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	item, err := GetByID(r)

	switch err {
	case nil:
		break
	case ErrParseItemID:
		http.Error(w, ErrParseItemID.Error(), http.StatusBadRequest)
		return
	case inventory.ErrItemNotFound:
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	default:
		http.Error(w, serverError, http.StatusInternalServerError)
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
		http.Error(w, serverError, http.StatusInternalServerError)
	}
}

func GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	all, err := ListItems(r)

	switch err {
	case nil:
		break
	default:
		http.Error(w, serverError, http.StatusInternalServerError)
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
		http.Error(w, serverError, http.StatusInternalServerError)
	}
}
