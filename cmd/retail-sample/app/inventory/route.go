package inventory

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/middleware"
)

func ConfigureRoutes(r *mux.Router, logger middleware.Logger, loggerFactory middleware.NewLoggerFunc, factory middleware.PersistenceProviderFactory) {
	items := webApp{
		logger: logger,
		wrapper: wrapper{
			Middleware: middleware.Middleware{
				NewLogger:                  loggerFactory,
				PersistenceProviderFactory: factory,
			},
		},
	}

	r.HandleFunc("/inventory", items.getAll).Methods(http.MethodGet)
	r.HandleFunc("/inventory/{itemID}", items.get).Methods(http.MethodGet)
	r.HandleFunc("/inventory/{itemID}", items.update).Methods(http.MethodPatch)
	r.HandleFunc("/inventory", items.create).Methods(http.MethodPost)
}
