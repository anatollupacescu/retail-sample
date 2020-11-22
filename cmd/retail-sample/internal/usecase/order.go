package usecase

import (
	"context"

	"github.com/anatollupacescu/retail-sample/domain/retail/order"
)

func NewOrder(ctx context.Context, orders order.Orders, db orderDB, log logger) Order {
	return Order{
		ctx:     ctx,
		orders:  orders,
		orderDB: db,
		logger:  log,
	}
}

type orderDB interface {
	Get(order.ID) (order.Order, error)
}

type Order struct {
	logger  logger
	orders  order.Orders
	orderDB orderDB
	ctx     context.Context
}

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

	newOrder, err := o.orderDB.Get(id)
	if err != nil {
		o.logger.Error("create order", "call domain to retrieve new order", err)
		return order.Order{}, err
	}

	o.logger.Info("create order", "success")

	return newOrder, nil
}
