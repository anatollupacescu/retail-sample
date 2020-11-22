package order

import (
	"github.com/anatollupacescu/retail-sample/domain/retail/order"
)

type PlaceOrderDTO struct {
	RecipeID, OrderQty int
}

func (o *Order) PlaceOrder(dto PlaceOrderDTO) (order.Order, error) {
	o.logger.Info("create order", "enter")

	id, err := o.orders.PlaceOrder(dto.RecipeID, dto.OrderQty)
	if err != nil {
		o.logger.Error("create order", "call domain layer", err)
		return order.Order{}, err
	}

	newOrder, err := o.orders.Get(id)
	if err != nil {
		o.logger.Error("create order", "call domain to retrieve new order", err)
		return order.Order{}, err
	}

	o.logger.Info("create order", "success")

	return newOrder, nil
}
