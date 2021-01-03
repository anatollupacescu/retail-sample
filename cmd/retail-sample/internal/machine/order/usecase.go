package order

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	pg "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/persistence/postgres"
	"github.com/anatollupacescu/retail-sample/domain/retail/order"
	"github.com/anatollupacescu/retail-sample/domain/retail/recipe"
	"github.com/anatollupacescu/retail-sample/domain/retail/stock"
)

type UseCase struct {
	ctx     context.Context
	logger  *zerolog.Logger
	orderDB *pg.OrderPgxStore
	orders  order.Orders
}

func New(ctx context.Context, t pg.TX) UseCase {
	logger := log.Ctx(ctx).With().Str("layer", "use case").Logger()

	orderDB := &pg.OrderPgxStore{DB: t.Tx}
	recipeDB := &pg.RecipePgxStore{DB: t.Tx}
	stockDB := &pg.StockPgxStore{DB: t.Tx}

	orders := order.Orders{
		DB: orderDB,
		Stock: &stock.Extractor{
			Recipes: recipeDB,
			Stock:   stockDB,
		},
		Recipes: &recipe.Validator{
			Recipes: recipeDB,
		},
	}

	return UseCase{
		ctx:     ctx,
		orders:  orders,
		orderDB: orderDB,
		logger:  &logger,
	}
}
