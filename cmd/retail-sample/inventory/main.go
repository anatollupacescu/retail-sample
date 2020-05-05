package inventory

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/types"
)

func ConfigureRoutes(r *mux.Router, loggerFactory types.LoggerFactory, factory types.PersistenceProviderFactory) {
	app := InventoryWebApp{
		Logger: loggerFactory(),
		Wrapper: InventoryWrapper{
			LoggerFactory:              loggerFactory,
			PersistenceProviderFactory: factory,
		},
	}
	r.HandleFunc("/inventory", app.GetAllInventoryItems).Methods(http.MethodGet)
	r.HandleFunc("/inventory/{itemID}", app.GetInventoryItem).Methods(http.MethodGet)
	r.HandleFunc("/inventory", app.CreateInventoryItem).Methods(http.MethodPost)
}
