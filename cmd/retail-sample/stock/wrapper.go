package stock

import (
	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/types"

	"github.com/anatollupacescu/retail-sample/internal/retail-domain/stock"
)

type wrapper struct {
	loggerFactory              types.LoggerFactory
	persistenceProviderFactory types.PersistenceProviderFactory

	provider types.PersistenceProvider
	logger   types.Logger
}

func (w *wrapper) exec(methodName string, f func() error) {
	logger := w.loggerFactory()

	logger.Log("msg", "enter", "method", methodName)
	defer logger.Log("msg", "exit", "method", methodName)

	w.logger = logger
	w.provider = w.persistenceProviderFactory.New()

	err := f()

	if err != nil {
		logger.Log("msg", "rollback")
		w.persistenceProviderFactory.Rollback(w.provider)
		return
	}

	logger.Log("msg", "commit")
	w.persistenceProviderFactory.Commit(w.provider)
}

func (w wrapper) quantity(id int) (qty int, err error) {
	w.exec("get stock quantity", func() error {
		s := w.provider.Stock()
		qty, err = s.Quantity(id)

		return err
	})

	return
}

func (w wrapper) currentStock() (currentStock []stock.StockPosition, err error) {
	w.exec("get current stock", func() error {
		s := w.provider.Stock()

		currentStock, err = s.CurrentStock()

		return err
	})

	return
}

func (w wrapper) provision(in []stock.ProvisionEntry) (updatedQtys map[int]int, err error) {
	w.exec("provision stock", func() error {
		s := w.provider.Stock()

		updatedQtys, err = s.Provision(in)

		return err
	})

	return
}

func (w wrapper) getProvisionLog() (pl []stock.ProvisionEntry, err error) {
	w.exec("get provision log", func() error {
		s := w.provider.Stock()

		pl, err = s.GetProvisionLog()

		return err
	})

	return
}
