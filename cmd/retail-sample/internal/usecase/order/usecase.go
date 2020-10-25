package order

import (
	"context"

	"github.com/anatollupacescu/retail-sample/domain/retail/order"
)

type logger interface {
	Error(string, string, error)
	Info(string, string)
}

func New(ctx context.Context, orders order.Orders, log logger) Order {
	return Order{
		ctx:    ctx,
		orders: orders,
		logger: log,
	}
}

type Order struct {
	logger logger
	orders order.Orders
	ctx    context.Context
}
