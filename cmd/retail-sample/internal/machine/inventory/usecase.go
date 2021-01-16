package inventory

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/middleware"
	pg "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/persistence/postgres"
	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
)

type UseCase struct {
	ctx         context.Context
	logger      *zerolog.Logger
	inventory   *inventory.Collection
	inventoryDB *pg.InventoryPgxStore
}

func New(ctx context.Context) (UseCase, error) {
	logger := log.Ctx(ctx).With().Str("domain", "inventory").Logger()

	tx, err := middleware.ExtractTransactionCtx(ctx)

	if err != nil {
		logger.Error().Str("action", "extract transaction").Err(err)
		return UseCase{}, err
	}

	db := &pg.InventoryPgxStore{DB: tx}

	uc := UseCase{
		ctx:         ctx,
		inventory:   &inventory.Collection{DB: db},
		inventoryDB: db,
		logger:      &logger,
	}

	return uc, nil
}
