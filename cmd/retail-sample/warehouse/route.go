package warehouse

import (
	"github.com/anatollupacescu/retail-sample/internal/warehouse"
	"github.com/gorilla/mux"
)

type App struct {
	stock warehouse.Stock
}

func ConfigureRoutes(r *mux.Router) {
	a := App{
		stock: warehouse.NewStock(),
	}
	r.HandleFunc("/inbound/config", a.ListTypes).Methods("GET")
	r.HandleFunc("/inbound/config", a.ConfigureType).Methods("POST")
	r.HandleFunc("/inbound", a.PlaceInbound).Methods("POST")
	r.HandleFunc("/inbound", a.ListInbound).Methods("GET")
	r.HandleFunc("/stock", a.ShowStock).Methods("GET")
	r.HandleFunc("/outbound", a.PlaceOutbound).Methods("POST")
	r.HandleFunc("/outbound/config", a.ConfigureOutbound).Methods("POST")
}
