package web

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/persistence"
	retail "github.com/anatollupacescu/retail-sample/internal/retail-sample"
	"github.com/pkg/errors"
)

func (a *WebApp) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields() // catch unwanted fields

	var requestBody struct {
		ID  *int `json:"id"` // pointer so we can test for field absence
		Qty *int `json:"qty"`
	}

	if err := d.Decode(&requestBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if requestBody.ID == nil || requestBody.Qty == nil {
		http.Error(w, "name or quantity not provided", http.StatusBadRequest)
		return
	}

	recipeID := *requestBody.ID
	orderQty := *requestBody.Qty

	entryID, err := a.App.PlaceOrder(recipeID, orderQty)

	if err != nil {
		switch errors.Cause(err) {
		case retail.BusinessErr:
			http.Error(w, err.Error(), http.StatusBadRequest)
			break
		case persistence.DBErr:
			log.Printf("request to place order: %v", err)
			fallthrough
		default:
			http.Error(w, "unexpected error", http.StatusInternalServerError)
		}
		return
	}

	type DescriptorEntity struct {
		ID       int `json:"id"`
		RecipeID int `json:"recipeID"`
		Qty      int `json:"qty"`
	}

	var response = struct {
		Data DescriptorEntity `json:"data"`
	}{
		Data: DescriptorEntity{
			ID:       int(entryID),
			RecipeID: recipeID,
			Qty:      orderQty,
		},
	}

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (a *WebApp) GetOrder(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
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

	switch err {
	//TODO
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
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
