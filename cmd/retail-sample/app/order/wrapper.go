package order

import (
	types "github.com/anatollupacescu/retail-sample/cmd/retail-sample/middleware"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/order"
)

type wrapper struct {
	types.Wrapper
}

func (w wrapper) create(id int, qty int) (orderID order.ID, err error) {
	return orderID, w.Exec("add new order", func(provider types.PersistenceProvider) error {
		o := provider.Orders()

		orderID, err = o.PlaceOrder(id, qty)

		return err
	})
}

func (w wrapper) get(id order.ID) (ordr order.Order, err error) {
	return ordr, w.Exec("get order by id", func(provider types.PersistenceProvider) error {
		o := provider.Orders()

		ordr, err = o.Get(id)

		return err
	})
}

func (w wrapper) getAll() (orders []order.Order, err error) {
	return orders, w.Exec("list past orders", func(provider types.PersistenceProvider) error {
		o := provider.Orders()

		orders, err = o.List()

		return err
	})
}
