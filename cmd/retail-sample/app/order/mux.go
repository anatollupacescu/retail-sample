package order

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/anatollupacescu/retail-sample/internal/retail-domain/inventory"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/order"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/stock"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/middleware"
)

type (
	webApp struct {
		logger  middleware.Logger
		wrapper wrapper
	}

	entity struct {
		ID       int `json:"id"`
		RecipeID int `json:"recipeID"`
		Qty      int `json:"qty"`
	}
)

var internalServerError = "internal server error"

func (a webApp) create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var requestBody struct {
		ID  int `json:"id"`
		Qty int `json:"qty"`
	}

	if err := d.Decode(&requestBody); err != nil {
		a.logger.Log("action", "decode request payload", "error", err, "method", "order.create")
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	recipeID := requestBody.ID
	orderQty := requestBody.Qty

	entryID, err := a.wrapper.create(recipeID, orderQty)

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

	var response = struct {
		Data entity `json:"data"`
	}{
		Data: entity{
			ID:       int(entryID),
			RecipeID: recipeID,
			Qty:      orderQty,
		},
	}

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		a.logger.Log("action", "encode response", "error", err, "method", "order.create")
		http.Error(w, internalServerError, http.StatusBadRequest)
	}
}

func (a webApp) get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	rid := vars["orderID"]

	id, err := strconv.Atoi(rid)

	if err != nil {
		a.logger.Log("action", "parse id", "error", err, "method", "order.get")
		http.Error(w, "could not parse id", http.StatusBadRequest)
		return
	}

	orderID := order.ID(id)

	ordr, err := a.wrapper.get(orderID)

	switch err {
	case nil:
		break
	case order.ErrOrderNotFound:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		a.logger.Log("action", "call application", "error", err, "method", "order.create")
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	var response = struct {
		Data entity `json:"data"`
	}{
		Data: entity{
			ID:       int(ordr.ID),
			RecipeID: ordr.RecipeID,
			Qty:      ordr.Qty,
		},
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		a.logger.Log("action", "encode response", "error", err, "method", "order.create")
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}
}

func (a webApp) getAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var response struct {
		Data []entity `json:"data"`
	}

	response.Data = make([]entity, 0)

	all, err := a.wrapper.getAll()

	if err != nil {
		http.Error(w, internalServerError, http.StatusBadRequest)
		return
	}

	for _, ordr := range all {
		e := entity{
			ID:       int(ordr.ID),
			RecipeID: ordr.RecipeID,
			Qty:      ordr.Qty,
		}

		response.Data = append(response.Data, e)
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		a.logger.Log("action", "encode response", "error", err, "method", "order.create")
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}
}
