package stock

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/types"
)

func ConfigureRoutes(r *mux.Router, loggerFactory types.LoggerFactory, factory types.PersistenceProviderFactory) {
	app := webApp{
		logger: loggerFactory(),
		wrapper: wrapper{
			loggerFactory:              loggerFactory,
			persistenceProviderFactory: factory,
		},
	}

	r.HandleFunc("/stock/provisionlog", app.GetProvisionLog).Methods(http.MethodGet)

	r.HandleFunc("/stock", app.GetStock).Methods(http.MethodGet)
	r.HandleFunc("/stock/{itemID}", app.GetStockPosition).Methods(http.MethodGet)
	r.HandleFunc("/stock", app.ProvisionStock).Methods(http.MethodPost)
}
