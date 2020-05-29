package stock

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/types"
)

func ConfigureRoutes(r *mux.Router, logger types.Logger, loggerFactory types.LoggerFactory, factory types.PersistenceProviderFactory) {
	stock := webApp{
		logger: logger,
		wrapper: wrapper{
			loggerFactory:              loggerFactory,
			persistenceProviderFactory: factory,
		},
	}

	r.HandleFunc("/stock/provisionlog", stock.getProvisionLog).Methods(http.MethodGet)

	r.HandleFunc("/stock", stock.getAll).Methods(http.MethodGet)
	r.HandleFunc("/stock/{itemID}", stock.get).Methods(http.MethodGet)
	r.HandleFunc("/stock", stock.update).Methods(http.MethodPost)
}
