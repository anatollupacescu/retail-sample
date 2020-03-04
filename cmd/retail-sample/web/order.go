package web

import (
	"encoding/json"
	"net/http"
	"time"
)

func (a *WebAdapter) PlaceOrder(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "name and quantity not provided", http.StatusBadRequest)
		return
	}

	if err := a.App.PlaceOrder(*requestBody.ID, *requestBody.Qty); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *WebAdapter) GetOrder(w http.ResponseWriter, r *http.Request) {
	panic("should return order entity")
}

func (a *WebAdapter) ListOrders(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	type entry struct {
		Date     time.Time `json:"date"`
		RecipeID int       `json:"recipeID"`
		Qty      int       `json:"qty"`
	}

	var response struct {
		Data []entry `json:"data"`
	}

	response.Data = make([]entry, 0)

	for _, o := range a.Orders.All() {
		e := entry{
			Date:     o.Date,
			RecipeID: o.RecipeID,
			Qty:      o.Qty,
		}

		response.Data = append(response.Data, e)
	}

	err := json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
