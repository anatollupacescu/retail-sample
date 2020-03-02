package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/anatollupacescu/retail-sample/internal/retail-sample/inventory"
)

const (
	ErrUnique = "ERR_UNIQUE"
	ErrNoName = "ERR_NO_NAME"
)

func (a *App) CreateInventoryItem(w http.ResponseWriter, r *http.Request) {
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
		createdID, err := a.inventory.Add(itemName)

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

func (a *App) GetInventoryItems(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type entry struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	var response struct {
		Data []entry `json:"data"`
	}

	response.Data = make([]entry, 0)

	for _, tp := range a.inventory.All() {
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

func (a *App) GetStock(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type entry struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Qty  int    `json:"qty"`
	}

	var response struct {
		Data []entry `json:"data"`
	}

	response.Data = make([]entry, 0)

	for _, position := range a.stock.CurrentState() {
		response.Data = append(response.Data, entry{
			ID:   position.ID,
			Name: position.Name,
			Qty:  position.Qty,
		})
	}

	err := json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (a *App) ProvisionStock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields() // catch unwanted fields

	var requestBody map[int]int

	if err := d.Decode(&requestBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	type DescriptorEntity map[int]int
	descriptors := make(DescriptorEntity, len(requestBody))

	for id, qty := range requestBody {
		newQty, err := a.stock.Provision(id, qty)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		descriptors[id] = newQty
	}

	var response = struct {
		Data DescriptorEntity `json:"data"`
	}{
		Data: descriptors,
	}

	switch len(descriptors) {
	case 0:
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusCreated)
	}

	err := json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (a *App) GetProvisionLog(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type entry struct {
		Time time.Time `json:"time"`
		ID   int       `json:"id"`
		Qty  int       `json:"qty"`
	}

	var response struct {
		Data []entry `json:"data"`
	}

	for _, in := range a.stock.ProvisionLog() {
		response.Data = append(response.Data, entry{
			Time: in.Time,
			ID:   int(in.ID),
			Qty:  in.Qty,
		})
	}

	err := json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
