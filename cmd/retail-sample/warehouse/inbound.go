package warehouse

import (
	"encoding/json"
	"fmt"
	"github.com/anatollupacescu/retail-sample/internal/warehouse"
	"log"
	"net/http"
)

func (a *App) ConfigureType(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields() // catch unwanted fields

	var types []string

	if err := d.Decode(&types); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(types) == 0 {
		http.Error(w, "at least one element expected", http.StatusBadRequest)
		return
	}

	for _, t := range types {
		if err := a.stock.ConfigureInboundType(t); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			if _, err = fmt.Fprintf(w, "could not add item type '%s': %v\n", t, err); err != nil {
				log.Fatal(err)
			}
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *App) ListTypes(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	e := json.NewEncoder(w)

	if err := e.Encode(a.stock.ItemTypes()); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (a *App) ShowStock(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	e := json.NewEncoder(w)

	var response = make(map[string]int)

	for _, itemType := range a.stock.ItemTypes() {
		qty, err := a.stock.Quantity(itemType)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response[itemType] = qty
	}

	if err := e.Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (a *App) Provision(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields() // catch unwanted fields

	var t map[string]int

	if err := d.Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(t) == 0 {
		http.Error(w, "nothing to provision", http.StatusBadRequest)
		return
	}

	for key, value := range t {
		item := warehouse.InboundItem{
			Type: warehouse.InboundType(key),
			Qty:  value,
		}

		_, err := a.stock.Provision(item)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusAccepted)
}
