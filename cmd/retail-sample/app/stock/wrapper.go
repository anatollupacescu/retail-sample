package stock

import (
	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/middleware"
	"github.com/anatollupacescu/retail-sample/domain/retail-sample/inventory"
	"github.com/anatollupacescu/retail-sample/domain/retail-sample/stock"
)

type wrapper struct {
	middleware.Middleware
}

func (w wrapper) quantity(id int) (sp stock.Position, err error) {
	return sp, w.Exec("get stock quantity", func(provider middleware.PersistenceProvider) error {
		i := provider.Inventory()

		var item inventory.Item
		item, err = i.Get(id)

		if err != nil {
			return err
		}

		s := provider.Stock()

		var qty int
		qty, err = s.Quantity(item.ID)

		if err != nil {
			return err
		}

		sp = stock.Position{
			ID:   item.ID,
			Name: item.Name,
			Qty:  qty,
		}

		return err
	})
}

func (w wrapper) currentStock() (currentStock []stock.Position, err error) {
	return currentStock, w.Exec("get current stock", func(provider middleware.PersistenceProvider) error {
		s := provider.Stock()

		currentStock, err = s.CurrentStock()

		return err
	})
}

func (w wrapper) provision(id, qty int) (newQty int, err error) {
	return newQty, w.Exec("provision stock", func(provider middleware.PersistenceProvider) error {
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

		newQty, err = s.Quantity(itemID)

		return err
	})
}

func (w wrapper) getProvisionLog() (pl []stock.ProvisionEntry, err error) {
	return pl, w.Exec("get provision log", func(provider middleware.PersistenceProvider) error {
		s := provider.Stock()

		pl, err = s.GetAllProvisions()

		return err
	})
}
