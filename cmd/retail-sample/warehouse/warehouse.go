package warehouse

import (
	"encoding/json"
	"fmt"
	"github.com/anatollupacescu/retail-sample/internal/warehouse"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type App struct {
	stock warehouse.Stock
}

func ConfigureRoutes(r *mux.Router) {
	a := App{}
	r.HandleFunc("/inventory", a.ListTypes).Methods("GET")
	r.HandleFunc("/inventory", a.ConfigureType).Methods("POST")
	r.HandleFunc("/inventory/{name}/{qty:[0-9]+}", a.Provision).Methods("POST")
	r.HandleFunc("/stock", a.ShowStock).Methods("GET")
	r.HandleFunc("/stock", a.NewOutbound).Methods("POST")
	r.HandleFunc("/stock/configure", a.ConfigureOutbound).Methods("POST")
}

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
	for _, itemType := range a.stock.ItemTypes() {
		qty, err := a.stock.Quantity(itemType)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := fmt.Fprint(w, "for item type '", itemType, "' got qty: ", qty, "\n"); err != nil {
			log.Fatal(err)
		}
	}
}

func (a *App) Provision(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)

	vars := mux.Vars(r)
	itemTypeName := vars["name"]
	qtyParam := vars["qty"]

	qty, err := strconv.Atoi(qtyParam)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	item := warehouse.InboundItem{
		Type: warehouse.InboundType(itemTypeName),
		Qty:  qty,
	}

	newQty, err := a.stock.Provision(item)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := fmt.Fprint(w, "for item type '", item.Type, "' got new qty: ", newQty, "\n"); err != nil {
		log.Fatal(err)
		return
	}
}
