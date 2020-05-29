package inventory

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/types"
)

func ConfigureRoutes(r *mux.Router, logger types.Logger, loggerFactory types.LoggerFactory, factory types.PersistenceProviderFactory) {
	items := webApp{
		logger: logger,
		wrapper: wrapper{
			loggerFactory:              loggerFactory,
			persistenceProviderFactory: factory,
		},
	}
	r.HandleFunc("/inventory", items.getAll).Methods(http.MethodGet)
	r.HandleFunc("/inventory/{itemID}", items.get).Methods(http.MethodGet)
	r.HandleFunc("/inventory/{itemID}", items.update).Methods(http.MethodPatch)
	r.HandleFunc("/inventory", items.create).Methods(http.MethodPost)
}
