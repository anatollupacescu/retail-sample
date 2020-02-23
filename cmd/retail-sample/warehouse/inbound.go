package warehouse

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/anatollupacescu/retail-sample/internal/retail-sample/inventory"
)

const (
	ErrUnique = "ERR_UNIQUE"
	ErrNoName = "ERR_NO_NAME"
)

func (a *App) CreateInventoryItem(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields() // catch unwanted fields

	var types []string

	if err := d.Decode(&types); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var createdID int
	for _, t := range types {
		var err error
		if createdID, err = a.stock.AddInventoryName(t); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			var msg string
			switch err {
			case inventory.ErrEmptyName:
				msg = ErrNoName
			case inventory.ErrDuplicateName:
				msg = ErrUnique
			default:
				http.Error(w, "Server error", http.StatusInternalServerError)
				return
			}
			if _, err = fmt.Fprint(w, msg); err != nil {
				log.Fatal(err)
			}
			return
		}
	}

	if createdID == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *App) GetInventoryItems(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type Record struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	var result struct {
		Data []Record `json:"data"`
	}

	result.Data = make([]Record, 0)

	for _, tp := range a.stock.InventoryItems() {
		result.Data = append(result.Data, Record{
			ID:   int(tp.ID),
			Name: string(tp.Name),
		})
	}

	e := json.NewEncoder(w)

	if err := e.Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (a *App) GetStock(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type itm struct {
		Name string `json:"name"`
		Qty  int    `json:"qty"`
	}

	var result struct {
		Data []itm `json:"data"`
	}

	result.Data = make([]itm, 0)

	for _, position := range a.stock.CurrentState() {
		result.Data = append(result.Data, itm{
			Name: position.Name,
			Qty:  position.Qty,
		})
	}

	e := json.NewEncoder(w)
	if err := e.Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (a *App) ProvisionStock(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields() // catch unwanted fields

	var t map[int]int

	if err := d.Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var (
		provisionID int
		err         error
	)

	for id, qty := range t {
		if provisionID, err = a.stock.Provision(id, qty); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if provisionID == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (a *App) GetProvisionLog(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	e := json.NewEncoder(w)

	type inbound struct {
		Time time.Time `json:"time"`
		ID   int       `json:"id"`
		Qty  int       `json:"qty"`
	}

	var t struct {
		Inbound []inbound `json:"inbounds"`
	}

	for _, in := range a.stock.ProvisionLog() {
		e := inbound{
			Time: in.Time,
			ID:   int(in.ID),
			Qty:  in.Qty,
		}
		t.Inbound = append(t.Inbound, e)
	}

	if err := e.Encode(t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
