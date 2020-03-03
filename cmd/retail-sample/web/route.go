package web

import (
	"net/http"

	"github.com/gorilla/mux"
)

func ConfigureRoutes(r *mux.Router) {
	webApp := newInMemoryApp()

	r.HandleFunc("/log/provision", webApp.GetProvisionLog).Methods(http.MethodGet)

	r.HandleFunc("/inventory", webApp.GetAllInventoryItems).Methods(http.MethodGet)
	r.HandleFunc("/inventory/{itemID}", webApp.GetInventoryItem).Methods(http.MethodGet)
	r.HandleFunc("/inventory", webApp.CreateInventoryItem).Methods(http.MethodPost)

	r.HandleFunc("/stock", webApp.GetStock).Methods(http.MethodGet)
	r.HandleFunc("/stock/{itemID}", webApp.GetStockPosition).Methods(http.MethodGet)
	r.HandleFunc("/stock", webApp.ProvisionStock).Methods(http.MethodPost)

	r.HandleFunc("/order", webApp.PlaceOrder).Methods(http.MethodPost)
	r.HandleFunc("/order/{orderID}", webApp.GetOrder).Methods(http.MethodPost)
	r.HandleFunc("/order", webApp.ListOrders).Methods(http.MethodGet)

	r.HandleFunc("/recipe", webApp.ListRecipes).Methods(http.MethodGet)
	r.HandleFunc("/recipe/{recipeID}", webApp.GetRecipe).Methods(http.MethodGet)
	r.HandleFunc("/recipe", webApp.CreateRecipe).Methods(http.MethodPost)
}
