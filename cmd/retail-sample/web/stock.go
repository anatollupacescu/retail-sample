package web

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

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

func (a *App) GetStockPosition(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	rid := vars["itemID"]

	itemID, err := strconv.Atoi(rid)

	if err != nil {
		http.Error(w, "invalid item id provided", http.StatusBadRequest)
		return
	}

	type entry struct {
		Qty int `json:"qty"`
	}

	qty := a.stock.Quantity(itemID)

	var response = struct {
		Data entry `json:"data"`
	}{
		Data: entry{
			Qty: qty,
		},
	}

	err = json.NewEncoder(w).Encode(response)

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

	type entry map[int]int
	data := make(entry, len(requestBody))

	for id, qty := range requestBody {
		newQty, err := a.stock.Provision(id, qty)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data[id] = newQty
	}

	var response = struct {
		Data entry `json:"data"`
	}{
		Data: data,
	}

	switch len(data) {
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
