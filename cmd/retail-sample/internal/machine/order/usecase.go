package order

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/middleware"
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

func New(ctx context.Context) (UseCase, error) {
	logger := log.Ctx(ctx).With().Str("domain", "order").Logger()

	tx, err := middleware.ExtractTransactionCtx(ctx)

	if err != nil {
		logger.Error().Str("action", "extract transaction").Err(err)
		return UseCase{}, err
	}

	orderDB := &pg.OrderPgxStore{DB: tx}
	recipeDB := &pg.RecipePgxStore{DB: tx}
	stockDB := &pg.StockPgxStore{DB: tx}

	orders := order.Orders{
		DB: orderDB,
		Stock: &stock.Extractor{
			Recipes: recipeDB,
			Stock:   stockDB,
		},
		RecipeValidator: &recipe.Validator{
			Recipes: recipeDB,
		},
	}

	uc := UseCase{
		ctx:     ctx,
		orders:  orders,
		orderDB: orderDB,
		logger:  &logger,
	}

	return uc, nil
}
