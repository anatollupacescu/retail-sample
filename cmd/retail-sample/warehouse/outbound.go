package warehouse

import (
	"encoding/json"
	"net/http"

	"github.com/anatollupacescu/retail-sample/internal/warehouse"
)

func (a *App) PlaceOutbound(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields() // catch unwanted fields

	t := struct {
		Name *string `json:"name"` // pointer so we can test for field absence
		Qty  *int    `json:"qty"`
	}{}

	if err := d.Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if t.Name == nil || t.Qty == nil {
		http.Error(w, "please provide valid name and quantity", http.StatusBadRequest)
		return
	}

	if err := a.stock.PlaceOutbound(*t.Name, *t.Qty); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a *App) ConfigureOutbound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields() // catch unwanted fields

	t := struct {
		Name  *string        `json:"name"` // pointer so we can test for field absence
		Items map[string]int `json:"items"`
	}{}

	if err := d.Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if t.Name == nil {
		http.Error(w, "name must not be empty", http.StatusBadRequest)
		return
	}

	var components []warehouse.OutboundItemComponent
	for itemType, qty := range t.Items {
		components = append(components, warehouse.OutboundItemComponent{
			ItemType: warehouse.InboundType(itemType),
			Qty:      qty,
		})
	}

	if err := a.stock.ConfigureOutbound(*t.Name, components); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
