package recipe

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/types"
)

func ConfigureRoutes(r *mux.Router, loggerFactory types.LoggerFactory, factory types.PersistenceProviderFactory) {
	app := webApp{
		logger: loggerFactory(),
		wrapper: wrapper{
			loggerFactory:              loggerFactory,
			persistenceProviderFactory: factory,
		},
	}

	r.HandleFunc("/recipe", app.ListRecipes).Methods(http.MethodGet)
	r.HandleFunc("/recipe/{recipeID}", app.GetRecipe).Methods(http.MethodGet)
	r.HandleFunc("/recipe", app.CreateRecipe).Methods(http.MethodPost)
}
