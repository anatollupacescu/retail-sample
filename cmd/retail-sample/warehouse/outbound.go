package warehouse

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/anatollupacescu/retail-sample/internal/retail-sample/recipe"
	"github.com/anatollupacescu/retail-sample/internal/retail-sample/warehouse"
)

func (a *App) CreateOutbound(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields() // catch unwanted fields

	t := struct {
		ID  *int `json:"name"` // pointer so we can test for field absence
		Qty *int `json:"qty"`
	}{}

	if err := d.Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if t.ID == nil || t.Qty == nil {
		http.Error(w, "please provide valid name and quantity", http.StatusBadRequest)
		return
	}

	if err := a.stock.PlaceOutbound(*t.ID, *t.Qty); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a *App) ListOutbound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	var (
		outbounds []warehouse.SoldItem
		err       error
	)

	if outbounds, err = a.stock.ListOutbound(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	type itm struct {
		Date time.Time `json:"date"`
		Name string    `json:"name"`
		Qty  int       `json:"qty"`
	}

	var result struct {
		Data []itm `json:"data"`
	}

	for _, o := range outbounds {
		result.Data = append(result.Data, itm{
			Date: o.Date,
			Name: o.Name,
			Qty:  o.Qty,
		})
	}

	e := json.NewEncoder(w)

	if err := e.Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (a *App) CreateOutboundConfig(w http.ResponseWriter, r *http.Request) {
	var t struct {
		Name  *string     `json:"name"` // pointer so we can test for field absence
		Items map[int]int `json:"items"`
	}

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields() // catch unwanted fields

	if err := d.Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if t.Name == nil {
		http.Error(w, "name must not be empty", http.StatusBadRequest)
		return
	}

	var ingredients []recipe.Ingredient
	for id, qty := range t.Items {
		ingredients = append(ingredients, recipe.Ingredient{
			ID:  id,
			Qty: qty,
		})
	}

	if err := a.stock.ConfigureOutbound(*t.Name, ingredients); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *App) ListOutboundConfig(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type itm struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}

	var result struct {
		Data []itm `json:"data"`
	}

	for _ = range /*on hold*/ []int{} {
		result.Data = append(result.Data, itm{
			// Name:  v.Name,
			// Count: len(v.Items),
		})
	}

	e := json.NewEncoder(w)

	if err := e.Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
