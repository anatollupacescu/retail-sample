package warehouse

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/anatollupacescu/retail-sample/internal/retail-sample/warehouse"
)

const ErrUnique = "ERR_UNIQUE"
const ErrNoName = "ERR_NO_NAME"

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
			var msg string
			switch err {
			case warehouse.ErrInboundNameNotProvided:
				msg = ErrNoName
			case warehouse.ErrInboundItemTypeAlreadyConfigured:
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

	w.WriteHeader(http.StatusCreated)
}

func (a *App) GetType(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	name := vars["name"]

	t := a.stock.GetType(name)

	var zeroItemConfig warehouse.ItemConfig

	if t == zeroItemConfig {
		http.Error(w, "item type not found", http.StatusNotFound)
		return
	}

	payload := struct {
		Disabled bool `json:"disabled"`
	}{
		Disabled: t.Disabled,
	}

	e := json.NewEncoder(w)

	if err := e.Encode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *App) DisableType(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	name := vars["name"]

	if err := a.stock.Disable(name); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (a *App) ListTypes(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type itm struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	}

	var result struct {
		Data []itm `json:"data"`
	}

	result.Data = make([]itm, 0)

	inboundTypes := a.stock.ItemTypes()
	for i, tp := range inboundTypes {
		result.Data = append(result.Data, itm{
			ID:   strconv.Itoa(i + 1),
			Type: tp,
		})
	}

	e := json.NewEncoder(w)

	if err := e.Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (a *App) ShowStock(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type itm struct {
		Name string `json:"name"`
		Qty  int    `json:"qty"`
	}

	var result struct {
		Data []itm `json:"data"`
	}

	result.Data = make([]itm, 0)

	for _, itemType := range a.stock.ItemTypes() {
		qty, err := a.stock.Quantity(itemType)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		result.Data = append(result.Data, itm{
			Name: itemType,
			Qty:  qty,
		})
	}

	e := json.NewEncoder(w)
	if err := e.Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (a *App) PlaceInbound(w http.ResponseWriter, r *http.Request) {
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
		item := warehouse.Item{
			ItemConfig: warehouse.ItemConfig{
				Type: key,
			},
			Qty: value,
		}

		_, err := a.stock.PlaceInbound(item)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusAccepted)
}

func (a *App) ListInbound(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	e := json.NewEncoder(w)

	type inbound struct {
		Name string `json:"name"`
		Qty  int    `json:"qty"`
	}

	var t struct {
		Inbound []inbound `json:"inbounds"`
	}

	for _, in := range a.stock.ListInbound() {
		e := inbound{
			Name: in.Type,
			Qty:  in.Qty,
		}
		t.Inbound = append(t.Inbound, e)
	}

	if err := e.Encode(t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
