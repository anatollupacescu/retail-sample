package stock

import (
	types "github.com/anatollupacescu/retail-sample/cmd/retail-sample/middleware"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/stock"
)

type wrapper struct {
	types.Wrapper
}

func (w wrapper) quantity(id int) (qty int, err error) {
	return qty, w.Exec("get stock quantity", func(provider types.PersistenceProvider) error {
		s := provider.Stock()

		qty, err = s.Quantity(id)

		return err
	})
}

func (w wrapper) currentStock() (currentStock []stock.StockPosition, err error) {
	return currentStock, w.Exec("get current stock", func(provider types.PersistenceProvider) error {
		s := provider.Stock()

		currentStock, err = s.CurrentStock()

		return err
	})
}

func (w wrapper) provision(id, qty int) (updatedQtys map[int]int, err error) {
	return updatedQtys, w.Exec("provision stock", func(provider types.PersistenceProvider) error {
		s := provider.Stock()

		var provisionID int

		provisionID, err = s.Provision(id, qty)

		if err != nil {
			return err
		}

		var logEntry stock.ProvisionEntry

		logEntry, err = s.GetProvision(provisionID)

		if err != nil {
			return err
		}

		itemID := logEntry.ID

		var stockQty int

		stockQty, err = s.Quantity(itemID)

		if err != nil {
			return err
		}

		updatedQtys = map[int]int{
			itemID: stockQty,
		}

		return err
	})
}

func (w wrapper) getProvisionLog() (pl []stock.ProvisionEntry, err error) {
	return pl, w.Exec("get provision log", func(provider types.PersistenceProvider) error {
		s := provider.Stock()

		pl, err = s.GetAllProvisions()

		return err
	})
}
