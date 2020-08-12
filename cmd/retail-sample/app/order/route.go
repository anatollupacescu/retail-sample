package order

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/middleware"
)

func ConfigureRoutes(r *mux.Router, logger middleware.Logger, loggerFactory middleware.NewLoggerFunc, factory middleware.PersistenceProviderFactory) {
	orders := webApp{
		logger: logger,
		wrapper: wrapper{
			Middleware: middleware.Middleware{
				NewLogger:                  loggerFactory,
				PersistenceProviderFactory: factory,
			},
		},
	}

	r.HandleFunc("/order", orders.create).Methods(http.MethodPost)
	r.HandleFunc("/order/{orderID}", orders.get).Methods(http.MethodGet)
	r.HandleFunc("/order", orders.getAll).Methods(http.MethodGet)
}
