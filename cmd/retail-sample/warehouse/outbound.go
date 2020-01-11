package warehouse

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/anatollupacescu/retail-sample/internal/retail-sample/warehouse"
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

func (a *App) ConfigureOutbound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields() // catch unwanted fields

	var t struct {
		Name  *string        `json:"name"` // pointer so we can test for field absence
		Items map[string]int `json:"items"`
	}

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
			ItemType: itemType,
			Qty:      qty,
		})
	}

	if err := a.stock.ConfigureOutbound(*t.Name, components); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
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

	for _, v := range a.stock.OutboundConfigurations() {
		result.Data = append(result.Data, itm{
			Name:  v.Name,
			Count: len(v.Items),
		})
	}

	e := json.NewEncoder(w)

	if err := e.Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
