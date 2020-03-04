package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/anatollupacescu/retail-sample/internal/retail-domain/inventory"
)

const (
	ErrUnique = "ERR_UNIQUE"
	ErrNoName = "ERR_NO_NAME"
)

func (a *WebAdapter) CreateInventoryItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields() // catch unwanted fields

	var names []string

	if err := d.Decode(&names); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	type DescriptorEntity map[inventory.Name]inventory.ID

	descriptors := make(DescriptorEntity, len(names))

	for _, name := range names {
		itemName := inventory.Name(name)
		createdID, err := a.Inventory.Add(itemName)

		if err != nil {
			var msg string
			switch err {
			case inventory.ErrEmptyName:
				msg = ErrNoName
			case inventory.ErrDuplicateName:
				msg = ErrUnique
			default:
				http.Error(w, "Server error", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusBadRequest)

			if _, err = fmt.Fprint(w, msg); err != nil {
				log.Fatal(err)
			}

			return
		}

		descriptors[itemName] = createdID
	}

	if len(descriptors) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var response = struct {
		Data DescriptorEntity `json:"data"`
	}{
		Data: descriptors,
	}

	w.WriteHeader(http.StatusCreated)

	err := json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (a *WebAdapter) GetAllInventoryItems(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type entry struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	var response struct {
		Data []entry `json:"data"`
	}

	response.Data = make([]entry, 0)

	for _, tp := range a.Inventory.All() {
		response.Data = append(response.Data, entry{
			ID:   int(tp.ID),
			Name: string(tp.Name),
		})
	}

	err := json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func isValidItemID(rid string) bool {
	if len(rid) == 0 {
		return false
	}

	if i, err := strconv.Atoi(rid); err != nil || i == 0 {
		return false
	}

	return true
}

func (a *WebAdapter) GetInventoryItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	rid := vars["itemID"]

	if !isValidItemID(rid) {
		http.Error(w, "invalid item id provided", http.StatusBadRequest)
		return
	}

	i, _ := strconv.Atoi(rid)

	inventoyID := inventory.ID(i)

	inventoryItem := a.Inventory.Get(inventoyID)

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

	err := json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
