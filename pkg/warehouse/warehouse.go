package warehouse

import (
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

func (a *App) ConfigureType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	itemTypeName := vars["name"]
	if err := a.stock.ConfigureInboundType(itemTypeName); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if _, err = fmt.Fprintf(w, "could not add item type '%s': %v\n", itemTypeName, err); err != nil {
			log.Fatal(err)
		}
	} else {
		w.WriteHeader(http.StatusCreated)
	}
}

func (a *App) ListTypes(w http.ResponseWriter, _ *http.Request) {
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := fmt.Fprint(w, "for item type '", item.Type, "' got new qty: ", newQty, "\n"); err != nil {
		log.Fatal(err)
	}
}
