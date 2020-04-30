package web

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/anatollupacescu/retail-sample/internal/retail-domain/inventory"
)

func (a *WebApp) CreateInventoryItem(w http.ResponseWriter, r *http.Request) {
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
	createdID, err := a.AddToInventory(itemName)

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

func (a *WebApp) GetAllInventoryItems(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	list, err := a.ListInventoryItems()

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

func (a *WebApp) GetInventoryItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	rid := vars["itemID"]

	i, _ := strconv.Atoi(rid)

	inventoyID := inventory.ID(i)

	inventoryItem, err := a.Inventory.Get(inventoyID)

	switch err {
	case nil:
		break
	case inventory.ErrInventoryItemNotFound:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		a.Logger.Log("action", "call application", "error", err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	type entry struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	var response = struct {
		Data entry `json:"data"`
	}{
		Data: entry{
			ID:   int(inventoryItem.ID),
			Name: string(inventoryItem.Name),
		},
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		a.Logger.Log("action", "encode response", "error", err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}
}
