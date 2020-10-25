package order

import (
	"github.com/anatollupacescu/retail-sample/domain/retail/order"
)

func (o *Order) Create(recipeID, qty int) (id order.ID, err error) {
	o.logger.Info("get all", "enter")

	id, err = o.orders.PlaceOrder(recipeID, qty)

	if err != nil {
		o.logger.Error("get all", "call domain layer", err)
		return
	}

	o.logger.Info("get all", "success")

	return
}
