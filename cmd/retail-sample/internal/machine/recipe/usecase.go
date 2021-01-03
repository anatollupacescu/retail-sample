package recipe

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/middleware"
	pg "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/persistence/postgres"
	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
	"github.com/anatollupacescu/retail-sample/domain/retail/recipe"
)

type UseCase struct {
	ctx      context.Context
	logger   *zerolog.Logger
	recipes  recipe.Recipes
	recipeDB *pg.RecipePgxStore
}

func New(ctx context.Context) (UseCase, error) {
	logger := log.Ctx(ctx).With().Str("layer", "use case").Logger()

	tx, err := middleware.ExtractTransactionCtx(ctx)

	if err != nil {
		return UseCase{}, err
	}

	recipeDB := &pg.RecipePgxStore{DB: tx}
	inventoryDB := &pg.InventoryPgxStore{DB: tx}

	recipes := recipe.Recipes{
		DB:        recipeDB,
		Validator: &inventory.Validator{Inventory: inventoryDB},
	}

	uc := UseCase{
		ctx:      ctx,
		recipes:  recipes,
		recipeDB: recipeDB,
		logger:   &logger,
	}

	return uc, nil
}
