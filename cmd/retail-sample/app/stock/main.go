package stock

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/middleware"
)

func ConfigureRoutes(r *mux.Router, logger middleware.Logger, loggerFactory middleware.LoggerFactory, factory middleware.PersistenceProviderFactory) {
	stock := webApp{
		logger: logger,
		wrapper: wrapper{
			Wrapper: middleware.Wrapper{
				LoggerFactory:              loggerFactory,
				PersistenceProviderFactory: factory,
			},
		},
	}

	r.HandleFunc("/stock/provisionlog", stock.getProvisionLog).Methods(http.MethodGet)

	r.HandleFunc("/stock", stock.getAll).Methods(http.MethodGet)
	r.HandleFunc("/stock/{itemID}", stock.get).Methods(http.MethodGet)
	r.HandleFunc("/stock", stock.update).Methods(http.MethodPost)
}
