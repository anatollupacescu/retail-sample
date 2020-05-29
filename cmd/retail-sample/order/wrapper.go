package order

import (
	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/types"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/order"
)

type wrapper struct {
	loggerFactory              types.LoggerFactory
	persistenceProviderFactory types.PersistenceProviderFactory

	provider types.PersistenceProvider
	logger   types.Logger
}

func (w *wrapper) exec(methodName string, f func(o order.Orders) error) {
	w.logger = w.loggerFactory()

	w.logger.Log("msg", "enter", "method", methodName)
	defer w.logger.Log("msg", "exit", "method", methodName)

	w.provider = w.persistenceProviderFactory.New()

	orders := w.provider.Orders()

	err := f(orders)

	if err != nil {
		w.logger.Log("msg", "rollback")
		w.persistenceProviderFactory.Rollback(w.provider)
		return
	}

	w.logger.Log("msg", "commit")
	w.persistenceProviderFactory.Commit(w.provider)
}

func (w wrapper) create(id int, qty int) (orderID order.ID, err error) {
	w.exec("add new order", func(o order.Orders) error {
		orderID, err = o.PlaceOrder(id, qty)

		if err != nil {
			w.logger.Log("error", err)
		}

		return err
	})

	return
}

func (w wrapper) get(id order.ID) (ordr order.Order, err error) {
	w.exec("get order by id", func(o order.Orders) error {
		ordr, err = o.Get(id)

		return err
	})

	return
}

func (w wrapper) getAll() (orders []order.Order, err error) {
	w.exec("list past orders", func(o order.Orders) error {
		orders, err = o.List()

		return err
	})

	return
}
