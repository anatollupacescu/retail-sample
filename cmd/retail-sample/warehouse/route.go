package warehouse

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/anatollupacescu/retail-sample/internal/retail-sample/warehouse"
)

type App struct {
	stock warehouse.Stock
}

func ConfigureRoutes(r *mux.Router) {
	a := App{
		stock: warehouse.NewInMemoryStock(),
	}

	r.HandleFunc("/log/provision", a.GetProvisionLog).Methods(http.MethodGet)
	r.HandleFunc("/log/order", a.GetOrderLog).Methods(http.MethodGet)

	r.HandleFunc("/inventory", a.GetInventoryItems).Methods(http.MethodGet)
	r.HandleFunc("/inventory", a.CreateInventoryItem).Methods(http.MethodPost)

	r.HandleFunc("/stock", a.GetStock).Methods(http.MethodGet)
	r.HandleFunc("/stock", a.ProvisionStock).Methods(http.MethodPost)

	r.HandleFunc("/order", a.PlaceOrder).Methods(http.MethodPost)

	r.HandleFunc("/recipe", a.ListRecipes).Methods(http.MethodGet)
	r.HandleFunc("/recipe", a.CreateRecipe).Methods(http.MethodPost)
}
