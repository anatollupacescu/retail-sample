package order

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/types"
)

func ConfigureRoutes(r *mux.Router, logger types.Logger, loggerFactory types.LoggerFactory, factory types.PersistenceProviderFactory) {
	order := orderWebApp{
		logger: logger,
		wrapper: wrapper{
			loggerFactory:              loggerFactory,
			persistenceProviderFactory: factory,
		},
	}
	r.HandleFunc("/order", order.create).Methods(http.MethodPost)
	r.HandleFunc("/order/{orderID}", order.get).Methods(http.MethodGet)
	r.HandleFunc("/order", order.getAll).Methods(http.MethodGet)
}
