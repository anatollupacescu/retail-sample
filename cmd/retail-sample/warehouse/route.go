package warehouse

import (
	"github.com/gorilla/mux"

	"github.com/anatollupacescu/retail-sample/internal/retail-sample/warehouse"
)

type App struct {
	stock warehouse.Stock
}

func ConfigureRoutes(r *mux.Router) {
	a := App{
		stock: warehouse.NewInMemoryStock(),
	}
	r.HandleFunc("/inbound/config", a.ListTypes).Methods("GET")
	r.HandleFunc("/inbound/config", a.ConfigureType).Methods("POST")
	r.HandleFunc("/inbound", a.PlaceInbound).Methods("POST")
	r.HandleFunc("/inbound", a.ListInbound).Methods("GET")
	r.HandleFunc("/stock", a.ShowStock).Methods("GET")
	r.HandleFunc("/outbound", a.PlaceOutbound).Methods("POST")
	r.HandleFunc("/outbound/config", a.ConfigureOutbound).Methods("POST")
	r.HandleFunc("/outbound/config", a.ListOutbound).Methods("GET")
}
