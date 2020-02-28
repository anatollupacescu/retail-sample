package web

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/anatollupacescu/retail-sample/internal/retail-sample/recipe"
)

func (a *App) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields() // catch unwanted fields

	t := struct {
		ID  *int `json:"id"` // pointer so we can test for field absence
		Qty *int `json:"qty"`
	}{}

	if err := d.Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if t.ID == nil || t.Qty == nil {
		http.Error(w, "name and quantity not provided", http.StatusBadRequest)
		return
	}

	if err := a.stock.PlaceOrder(*t.ID, *t.Qty); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a *App) ListOrders(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	type itm struct {
		Date time.Time `json:"date"`
		Name string    `json:"name"`
		Qty  int       `json:"qty"`
	}

	var result struct {
		Data []itm `json:"data"`
	}

	for _, o := range a.stock.OrderLog() {
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

func (a *App) CreateRecipe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields() // catch unwanted fields

	var t struct {
		Name  *string     `json:"name"` // pointer so we can test for field absence
		Items map[int]int `json:"items"`
	}

	if err := d.Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if t.Name == nil {
		http.Error(w, "name can not be empty", http.StatusBadRequest)
		return
	}

	var ingredients []recipe.Ingredient
	for id, qty := range t.Items {
		ingredients = append(ingredients, recipe.Ingredient{
			ID:  id,
			Qty: qty,
		})
	}

	var recipeName = recipe.Name(*t.Name)

	recipeID, err := a.recipe.Add(recipeName, ingredients)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	type DescriptorEntity map[recipe.Name]recipe.ID

	var result = struct {
		Data DescriptorEntity `json:"data"`
	}{
		Data: DescriptorEntity{
			recipeName: recipeID,
		},
	}

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(result)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (a *App) ListRecipes(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var result struct {
		Data []string `json:"data"`
	}

	result.Data = make([]string, 0)

	for _, name := range a.recipe.Names() {
		result.Data = append(result.Data, string(name))
	}

	e := json.NewEncoder(w)
	if err := e.Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
