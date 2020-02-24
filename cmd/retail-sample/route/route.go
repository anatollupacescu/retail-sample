package route

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/anatollupacescu/retail-sample/internal/retail-sample/inventory"
	"github.com/anatollupacescu/retail-sample/internal/retail-sample/recipe"
	"github.com/anatollupacescu/retail-sample/internal/retail-sample/warehouse"
)

type App struct {
	inventory inventory.Inventory
	recipe    recipe.Book
	stock     warehouse.Stock
}

func ConfigureRoutes(r *mux.Router) {
	inventryStore := inventory.NewInMemoryStore()
	inventory := inventory.Inventory{Store: &inventryStore}

	recipeStore := recipe.NewInMemoryStore()
	recipeBook := recipe.Book{Store: &recipeStore, Inventory: &inventory}

	webApp := App{
		inventory: inventory,
		recipe:    recipeBook,
		stock: warehouse.Stock{
			Inventory:   inventory,
			RecipeBook:  recipeBook,
			InboundLog:  make(warehouse.InMemoryInboundLog),
			OutboundLog: make(warehouse.InMemoryOutboundLog),
			Data:        make(map[int]int),
		},
	}

	r.HandleFunc("/log/provision", webApp.GetProvisionLog).Methods(http.MethodGet)
	r.HandleFunc("/log/order", webApp.ListOrders).Methods(http.MethodGet)

	r.HandleFunc("/inventory", webApp.GetInventoryItems).Methods(http.MethodGet)
	r.HandleFunc("/inventory", webApp.CreateInventoryItem).Methods(http.MethodPost)

	r.HandleFunc("/stock", webApp.GetStock).Methods(http.MethodGet)
	r.HandleFunc("/stock", webApp.ProvisionStock).Methods(http.MethodPost)

	r.HandleFunc("/order", webApp.PlaceOrder).Methods(http.MethodPost)

	r.HandleFunc("/recipe", webApp.ListRecipes).Methods(http.MethodGet)
	r.HandleFunc("/recipe", webApp.CreateRecipe).Methods(http.MethodPost)
}
