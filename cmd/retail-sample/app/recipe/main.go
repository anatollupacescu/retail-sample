package recipe

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/middleware"
)

func ConfigureRoutes(r *mux.Router, logger middleware.Logger, loggerFactory middleware.LoggerFactory, factory middleware.PersistenceProviderFactory) {
	recipes := webApp{
		logger: logger,
		wrapper: wrapper{
			Wrapper: middleware.Wrapper{
				LoggerFactory:              loggerFactory,
				PersistenceProviderFactory: factory,
			},
		},
	}

	r.HandleFunc("/recipe", recipes.getAll).Methods(http.MethodGet)
	r.HandleFunc("/recipe/{recipeID}", recipes.get).Methods(http.MethodGet)
	r.HandleFunc("/recipe/{recipeID}", recipes.update).Methods(http.MethodPatch)
	r.HandleFunc("/recipe", recipes.create).Methods(http.MethodPost)
}
