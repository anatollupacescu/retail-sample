package stock

import (
	types "github.com/anatollupacescu/retail-sample/cmd/retail-sample/middleware"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/stock"
)

type wrapper struct {
	types.Wrapper
}

func (w wrapper) quantity(id int) (qty int, err error) {
	w.Exec("get stock quantity", func(provider types.PersistenceProvider) error {
		s := provider.Stock()

		qty, err = s.Quantity(id)

		return err
	})

	return
}

func (w wrapper) currentStock() (currentStock []stock.StockPosition, err error) {
	w.Exec("get current stock", func(provider types.PersistenceProvider) error {
		s := provider.Stock()

		currentStock, err = s.CurrentStock()

		return err
	})

	return
}

func (w wrapper) provision(in []stock.ProvisionEntry) (updatedQtys map[int]int, err error) {
	w.Exec("provision stock", func(provider types.PersistenceProvider) error {
		s := provider.Stock()

		updatedQtys, err = s.Provision(in)

		return err
	})

	return
}

func (w wrapper) getProvisionLog() (pl []stock.ProvisionEntry, err error) {
	w.Exec("get provision log", func(provider types.PersistenceProvider) error {
		s := provider.Stock()

		pl, err = s.GetProvisionLog()

		return err
	})

	return
}
