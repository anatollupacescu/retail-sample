package warehouse

import (
	"fmt"
	"github.com/anatollupacescu/retail-sample/internal/warehouse"
	"net/http"
)

type App struct {
	stock warehouse.Stock
}

func (a *App) InventoryHandler(w http.ResponseWriter, _ *http.Request)  {
	w.Header().Set("Content-type", "application/json")

	for _, itemType := range a.stock.ItemTypes() {
		line := fmt.Sprint("for item type '", itemType, " got qty: ", )
		w.Write([]byte(line))
	}

	w.WriteHeader(http.StatusOK)
}
