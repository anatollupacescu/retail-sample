package stock

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/middleware"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/inventory"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/stock"
)

type (
	webApp struct {
		logger  middleware.Logger
		wrapper wrapper
	}
	entity struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Qty  int    `json:"qty"`
	}
)

var internalServerError = "internal server error"

func (a *webApp) getAll(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	entries, err := a.wrapper.currentStock()

	if err != nil {
		http.Error(w, internalServerError, http.StatusBadRequest)
		return
	}

	var response struct {
		Data []entity `json:"data"`
	}

	response.Data = make([]entity, 0)

	for _, position := range entries {
		response.Data = append(response.Data, entity{
			ID:   position.ID,
			Name: position.Name,
			Qty:  position.Qty,
		})
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		a.logger.Log("action", "encode response", "error", err, "method", "stock.getAll")
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}
}

func (a *webApp) get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	rid := vars["itemID"]

	itemID, err := strconv.Atoi(rid)

	if err != nil {
		a.logger.Log("action", "parse id", "error", err, "method", "stock.get")
		http.Error(w, "could not parse id", http.StatusBadRequest)
		return
	}

	qty, err := a.wrapper.quantity(itemID)

	switch err {
	case nil:
		break
	case stock.ErrItemNotFound:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	var response = struct {
		Data entity `json:"data"`
	}{
		Data: entity{
			ID:  itemID,
			Qty: qty,
		},
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		a.logger.Log("action", "encode response", "error", err, "method", "stock.get")
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}
}

func (a *webApp) update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var requestBody map[int]int

	if err := d.Decode(&requestBody); err != nil {
		a.logger.Log("action", "decode request", "error", err, "method", "stock.update")
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	entries := make([]stock.ProvisionEntry, 0)

	for id, qty := range requestBody {
		entries = append(entries, stock.ProvisionEntry{
			ID:  id,
			Qty: qty,
		})
	}

	data, err := a.wrapper.provision(entries)

	switch err {
	case nil:
		break
	case inventory.ErrItemNotFound:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
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
		a.logger.Log("action", "encode response", "error", err, "method", "stock.update")
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}
}

func (a *webApp) getProvisionLog(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	provisionLog, err := a.wrapper.getProvisionLog()

	if err != nil {
		a.logger.Log("action", "call application", "error", err, "method", "stock.provisionlog")
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
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

	for _, in := range provisionLog {
		response.Data = append(response.Data, provisionLogEntity{
			ID:  int(in.ID),
			Qty: in.Qty,
		})
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		a.logger.Log("action", "encode response", "error", err, "method", "stock.provisionlog")
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}
}
