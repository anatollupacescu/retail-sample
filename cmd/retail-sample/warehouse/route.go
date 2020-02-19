package warehouse

import (
	"net/http"

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

	r.HandleFunc("/inbound", a.PlaceInbound).Methods(http.MethodPost)
	r.HandleFunc("/inbound", a.ListInbound).Methods(http.MethodGet)

	r.HandleFunc("/inbound/config", a.ListInboundConfig).Methods(http.MethodGet)
	r.HandleFunc("/inbound/config/{name}", a.GetInboundConfig).Methods(http.MethodGet)
	r.HandleFunc("/inbound/config", a.CreateInboundConfig).Methods(http.MethodPost)
	r.HandleFunc("/inbound/config/{name}/disable", a.DisableInboundConfig).Methods(http.MethodPatch)

	r.HandleFunc("/stock", a.ShowStock).Methods(http.MethodGet)

	r.HandleFunc("/outbound", a.CreateOutbound).Methods(http.MethodPost)
	r.HandleFunc("/outbound", a.ListOutbound).Methods(http.MethodGet)
	r.HandleFunc("/outbound/config", a.CreateOutboundConfig).Methods(http.MethodPost)
	r.HandleFunc("/outbound/config", a.ListOutboundConfig).Methods(http.MethodGet)
}
