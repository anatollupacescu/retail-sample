package order

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	pg "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/persistence/postgres"
	domain "github.com/anatollupacescu/retail-sample/domain/retail/order"
)

func New(ctx context.Context, t pg.TX) UseCase {
	logger := log.Ctx(ctx).With().Str("layer", "usecase").Logger()

	orderDB := &pg.OrderPgxStore{DB: t.Tx}
	recipeDB := &pg.RecipePgxStore{DB: t.Tx}
	stockDB := &pg.StockPgxStore{DB: t.Tx}

	orders := domain.Orders{
		DB:      orderDB,
		Stock:   &adapter{stock: stockDB},
		Recipes: recipeDB,
	}

	return UseCase{
		ctx:     ctx,
		orders:  orders,
		orderDB: orderDB,
		logger:  &logger,
	}
}

type UseCase struct {
	ctx     context.Context
	logger  *zerolog.Logger
	orderDB *pg.OrderPgxStore
	orders  domain.Orders
}

type PlaceOrderDTO struct {
	RecipeID, OrderQty int
}

func (o *UseCase) PlaceOrder(dto PlaceOrderDTO) (domain.OrderDTO, error) {
	id, err := o.orders.Add(dto.RecipeID, dto.OrderQty)
	if err != nil {
		o.logger.Error().Err(err).Msg("call domain layer")
		return domain.OrderDTO{}, err
	}

	newOrder, err := o.orderDB.Get(id)
	if err != nil {
		o.logger.Error().Err(err).Msg("retrieve new order")
		return domain.OrderDTO{}, err
	}

	return newOrder, nil
}
