package order

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/types"
)

func ConfigureRoutes(r *mux.Router, loggerFactory types.LoggerFactory, factory types.PersistenceProviderFactory) {
	app := orderWebApp{
		logger: loggerFactory(),
		wrapper: wrapper{
			loggerFactory:              loggerFactory,
			persistenceProviderFactory: factory,
		},
	}
	r.HandleFunc("/order", app.PlaceOrder).Methods(http.MethodPost)
	r.HandleFunc("/order/{orderID}", app.GetOrder).Methods(http.MethodPost)
	r.HandleFunc("/order", app.ListOrders).Methods(http.MethodGet)
}
