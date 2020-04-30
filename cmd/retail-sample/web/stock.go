package web

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/anatollupacescu/retail-sample/internal/retail-domain/inventory"
	retail "github.com/anatollupacescu/retail-sample/internal/retail-sample"
	"github.com/anatollupacescu/retail-sample/internal/retail-sample/stock"
	"github.com/gorilla/mux"
)

func (a *WebApp) GetStock(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	stockData, err := a.CurrentStock()

	if err != nil {
		a.Logger.Log("action", "call application", "error", err)
		http.Error(w, internalServerError, http.StatusBadRequest)
	}

	type entry struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Qty  int    `json:"qty"`
	}

	var response struct {
		Data []entry `json:"data"`
	}

	response.Data = make([]entry, 0)

	for _, position := range stockData {
		response.Data = append(response.Data, entry{
			ID:   position.ID,
			Name: position.Name,
			Qty:  position.Qty,
		})
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		a.Logger.Log("action", "encode response", "error", err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}
}

func (a *WebApp) GetStockPosition(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	rid := vars["itemID"]

	itemID, _ := strconv.Atoi(rid)

	qty, err := a.Quantity(itemID)

	switch err {
	case nil:
		break
	case stock.ErrItemNotFound:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		a.Logger.Log("action", "call application", "error", err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	type entry struct {
		Qty int `json:"qty"`
	}

	var response = struct {
		Data entry `json:"data"`
	}{
		Data: entry{
			Qty: qty,
		},
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		a.Logger.Log("action", "encode response", "error", err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}
}

func (a *WebApp) ProvisionStock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var requestBody map[int]int

	if err := d.Decode(&requestBody); err != nil {
		a.Logger.Log("action", "decode request", "error", err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	entries := make([]retail.ProvisionEntry, 0)

	for id, qty := range requestBody {
		entries = append(entries, retail.ProvisionEntry{
			ID:  id,
			Qty: qty,
		})
	}

	data, err := a.Provision(entries)

	switch err {
	case nil:
		break
	case inventory.ErrInventoryItemNotFound:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		a.Logger.Log("action", "call application", "error", err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}

	var response = struct {
		Data map[int]int `json:"data"`
	}{
		Data: data,
	}

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		a.Logger.Log("action", "encode response", "error", err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}
}

func (a *WebApp) GetProvisionLog(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	provisionLog, err := a.App.GetProvisionLog()

	if err != nil {
		a.Logger.Log("action", "call application", "error", err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	type entry struct {
		Time time.Time `json:"time"`
		ID   int       `json:"id"`
		Qty  int       `json:"qty"`
	}

	var response struct {
		Data []entry `json:"data"`
	}

	response.Data = make([]entry, 0)

	for _, in := range provisionLog {
		response.Data = append(response.Data, entry{
			ID:  int(in.ID),
			Qty: in.Qty,
		})
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		a.Logger.Log("action", "encode response", "error", err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}
}
