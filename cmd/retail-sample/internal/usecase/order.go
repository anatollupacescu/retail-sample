package usecase

import (
	"context"

	"github.com/anatollupacescu/retail-sample/domain/retail/order"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type orderDB interface {
	Get(order.ID) (order.Order, error)
}

func NewOrder(ctx context.Context, orders order.Orders, db orderDB) Order {
	logger := log.Ctx(ctx).With().Str("layer", "usecase").Logger()

	return Order{
		ctx:     ctx,
		orders:  orders,
		orderDB: db,
		logger:  &logger,
	}
}

type Order struct {
	logger  *zerolog.Logger
	orders  order.Orders
	orderDB orderDB
	ctx     context.Context
}

type PlaceOrderDTO struct {
	RecipeID, OrderQty int
}

func (o *Order) PlaceOrder(dto PlaceOrderDTO) (order.Order, error) {
	id, err := o.orders.Add(dto.RecipeID, dto.OrderQty)
	if err != nil {
		o.logger.Error().Err(err).Msg("call domain layer")
		return order.Order{}, err
	}

	newOrder, err := o.orderDB.Get(id)
	if err != nil {
		o.logger.Error().Err(err).Msg("retrieve new order")
		return order.Order{}, err
	}

	return newOrder, nil
}
