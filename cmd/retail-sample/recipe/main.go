package recipe

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/types"
)

func ConfigureRoutes(r *mux.Router, logger types.Logger, loggerFactory types.LoggerFactory, factory types.PersistenceProviderFactory) {
	recipes := webApp{
		logger: logger,
		wrapper: wrapper{
			loggerFactory:              loggerFactory,
			persistenceProviderFactory: factory,
		},
	}

	r.HandleFunc("/recipe", recipes.getAll).Methods(http.MethodGet)
	r.HandleFunc("/recipe/{recipeID}", recipes.get).Methods(http.MethodGet)
	r.HandleFunc("/recipe/{recipeID}", recipes.update).Methods(http.MethodPatch)
	r.HandleFunc("/recipe", recipes.create).Methods(http.MethodPost)
}
