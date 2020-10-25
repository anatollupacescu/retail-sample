package order

import (
	"github.com/anatollupacescu/retail-sample/domain/retail/order"
)

func (o *Order) GetAll() (orders []order.Order, err error) {
	o.logger.Info("get all", "enter")

	orders, err = o.orders.List()

	if err != nil {
		o.logger.Error("get all", "call domain layer", err)

		return
	}

	o.logger.Info("get all", "success")

	return
}
