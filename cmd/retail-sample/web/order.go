package web

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/anatollupacescu/retail-sample/internal/retail-domain/inventory"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/order"
	"github.com/anatollupacescu/retail-sample/internal/retail-sample/stock"
	"github.com/gorilla/mux"
)

func (a *WebApp) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var requestBody struct {
		ID  int `json:"id"`
		Qty int `json:"qty"`
	}

	if err := d.Decode(&requestBody); err != nil {
		a.Logger.Log("action", "decode request payload", "error", err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	recipeID := requestBody.ID
	orderQty := requestBody.Qty

	entryID, err := a.App.PlaceOrder(recipeID, orderQty)

	switch err {
	case nil:
		break
	case stock.ErrNotEnoughStock:
		fallthrough
	case inventory.ErrDuplicateName:
		fallthrough
	case inventory.ErrEmptyName:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	type descriptorEntity struct {
		ID       int `json:"id"`
		RecipeID int `json:"recipeID"`
		Qty      int `json:"qty"`
	}

	var response = struct {
		Data descriptorEntity `json:"data"`
	}{
		Data: descriptorEntity{
			ID:       int(entryID),
			RecipeID: recipeID,
			Qty:      orderQty,
		},
	}

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		a.Logger.Log("action", "encode response", "error", err)
		http.Error(w, internalServerError, http.StatusBadRequest)
	}
}

func (a *WebApp) GetOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	rid := vars["orderID"]

	id, _ := strconv.Atoi(rid)

	orderID := order.ID(id)

	ordr, err := a.Orders.Get(orderID)

	switch err {
	case nil:
		break
	case order.ErrOrderNotFound:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		a.Logger.Log("action", "call application", "error", err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	type entry struct {
		RecipeID int `json:"recipeID"`
		Qty      int `json:"qty"`
	}

	var response = struct {
		Data entry `json:"data"`
	}{
		Data: entry{
			RecipeID: ordr.RecipeID,
			Qty:      ordr.Qty,
		},
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		a.Logger.Log("action", "encode response", "error", err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}
}

func (a *WebApp) ListOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	type entry struct {
		ID       int `json:"id"`
		RecipeID int `json:"recipeID"`
		Qty      int `json:"qty"`
	}

	var response struct {
		Data []entry `json:"data"`
	}

	response.Data = make([]entry, 0)

	list, err := a.Orders.List()

	if err != nil {
		a.Logger.Log("action", "call application", "error", err)
		http.Error(w, internalServerError, http.StatusBadRequest)
	}

	for _, o := range list {
		e := entry{
			ID:       int(o.ID),
			RecipeID: o.RecipeID,
			Qty:      o.Qty,
		}

		response.Data = append(response.Data, e)
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		a.Logger.Log("action", "encode response", "error", err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}
}
