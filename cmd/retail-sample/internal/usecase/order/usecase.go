package order

import (
	"context"

	"github.com/anatollupacescu/retail-sample/domain/retail/order"
)

type logger interface {
	Error(string, string, error)
	Info(string, string)
}

func New(ctx context.Context, orders order.Orders, db orderDB, log logger) Order {
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
