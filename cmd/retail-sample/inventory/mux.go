package inventory

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/types"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/inventory"
)

type InventoryWebApp struct {
	Logger  types.Logger
	Wrapper InventoryWrapper
}

var internalServerError = "internal server error"

func (a *InventoryWebApp) UpdateItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	rid := vars["itemID"]

	id, _ := strconv.Atoi(rid)

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	type payload struct {
		Enabled bool `json:"enabled"`
	}

	var requestPayload payload

	if err := d.Decode(&requestPayload); err != nil {
		a.Logger.Log("action", "decode request payload", "error", err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	item, err := a.Wrapper.ChangeItemStatus(id, requestPayload.Enabled)

	switch err {
	case nil:
		break
	case inventory.ErrItemNotFound:
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	default:
		a.Logger.Log("action", "call application", "error", err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	type DescriptorEntity struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		Enabled bool   `json:"enabled"`
	}

	var response = struct {
		Data DescriptorEntity `json:"data"`
	}{
		Data: DescriptorEntity{
			ID:      id,
			Name:    item.Name,
			Enabled: item.Enabled,
		},
	}

	w.WriteHeader(http.StatusAccepted)

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		a.Logger.Log("action", "encode response", "error", err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}
}

func (a *InventoryWebApp) CreateInventoryItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	type payload struct {
		Name string `json:"name"`
	}

	var requestPayload payload

	if err := d.Decode(&requestPayload); err != nil {
		a.Logger.Log("action", "decode request payload", "error", err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	itemName := requestPayload.Name
	createdID, err := a.Wrapper.AddToInventory(itemName)

	switch err {
	case nil:
		break
	case inventory.ErrEmptyName:
		fallthrough
	case inventory.ErrDuplicateName:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		a.Logger.Log("action", "call application", "error", err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	type DescriptorEntity struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	var response = struct {
		Data DescriptorEntity `json:"data"`
	}{
		Data: DescriptorEntity{
			ID:   int(createdID),
			Name: requestPayload.Name,
		},
	}

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		a.Logger.Log("action", "encode response", "error", err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}
}

func (a *InventoryWebApp) GetAllInventoryItems(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	list, err := a.Wrapper.ListInventoryItems()

	if err != nil {
		a.Logger.Log("action", "call application", "error", err)
		http.Error(w, internalServerError, http.StatusBadRequest)
	}

	type entry struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	var response struct {
		Data []entry `json:"data"`
	}

	response.Data = make([]entry, 0)

	for _, tp := range list {
		response.Data = append(response.Data, entry{
			ID:   int(tp.ID),
			Name: string(tp.Name),
		})
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		a.Logger.Log("action", "encode response", "error", err)
		http.Error(w, internalServerError, http.StatusBadRequest)
	}
}

func (a *InventoryWebApp) GetInventoryItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	rid := vars["itemID"]

	id, _ := strconv.Atoi(rid)

	inventoryItem, err := a.Wrapper.GetInventoryItem(id)

	switch err {
	case nil:
		break
	case inventory.ErrItemNotFound:
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	default:
		a.Logger.Log("action", "call application", "error", err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	type entry struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		Enabled bool   `json:"enabled"`
	}

	var response = struct {
		Data entry `json:"data"`
	}{
		Data: entry{
			ID:      int(inventoryItem.ID),
			Name:    string(inventoryItem.Name),
			Enabled: inventoryItem.Enabled,
		},
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		a.Logger.Log("action", "encode response", "error", err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}
}
