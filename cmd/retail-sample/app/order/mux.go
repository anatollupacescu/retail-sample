package order

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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

	singleResponse struct {
		Data entity `json:"data"`
	}
)

func internalServerErrorMsg() string {
	return http.StatusText(http.StatusInternalServerError)
}

type createPayload struct {
	ID  int `json:"id"`
	Qty int `json:"qty"`
}

func (a webApp) create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var payload createPayload

	if err := d.Decode(&payload); err != nil {
		a.logger.Log("action", "decode request payload", "error", err, "method", "order.create")
		http.Error(w, internalServerErrorMsg(), http.StatusInternalServerError)

		return
	}

	recipeID := payload.ID
	orderQty := payload.Qty

	entryID, err := a.wrapper.create(recipeID, orderQty)

	switch err {
	case nil:
		break
	case order.ErrInvalidRecipe,
		order.ErrInvalidQuantity,
		stock.ErrNotEnoughStock,
		inventory.ErrDuplicateName,
		inventory.ErrEmptyName:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		http.Error(w, internalServerErrorMsg(), http.StatusInternalServerError)
		return
	}

	var response = singleResponse{
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
		http.Error(w, internalServerErrorMsg(), http.StatusBadRequest)
	}
}

func Create(recipeID, qty int, post func(io.Reader) (*http.Response, error)) (o order.Order, err error) {
	payload := createPayload{
		ID:  recipeID,
		Qty: qty,
	}

	var data []byte

	data, err = json.Marshal(payload)
	if err != nil {
		return o, err
	}

	body := bytes.NewReader(data)

	response, err := post(body)

	if err != nil {
		return o, err
	}

	if response.StatusCode != http.StatusCreated {
		return o, fmt.Errorf("unexpected status code: %v", response.StatusCode)
	}

	defer func() {
		_ = response.Body.Close()
	}()

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return o, err
	}

	var responseData singleResponse

	if err = json.Unmarshal(respBody, &responseData); err != nil {
		return o, err
	}

	return order.Order{
		ID: order.ID(responseData.Data.ID),
		Entry: order.Entry{
			Qty:      responseData.Data.Qty,
			RecipeID: responseData.Data.RecipeID,
		},
	}, err
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
		http.Error(w, internalServerErrorMsg(), http.StatusInternalServerError)

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
		http.Error(w, internalServerErrorMsg(), http.StatusInternalServerError)
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
		http.Error(w, internalServerErrorMsg(), http.StatusBadRequest)
		return
	}

	for i := range all {
		currentOrder := all[i]

		e := entity{
			ID:       int(currentOrder.ID),
			RecipeID: currentOrder.RecipeID,
			Qty:      currentOrder.Qty,
		}

		response.Data = append(response.Data, e)
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		a.logger.Log("action", "encode response", "error", err, "method", "order.create")
		http.Error(w, internalServerErrorMsg(), http.StatusInternalServerError)
	}
}
