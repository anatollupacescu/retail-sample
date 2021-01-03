package stock

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/middleware"
	pg "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/persistence/postgres"
)

type UseCase struct {
	logger      *zerolog.Logger
	stockDB     *pg.StockPgxStore
	inventoryDB *pg.InventoryPgxStore
	logDB       *pg.PgxProvisionLog
	ctx         context.Context
}

func New(ctx context.Context) (UseCase, error) {
	logger := log.Ctx(ctx).With().Str("layer", "use case").Logger()

	tx, err := middleware.ExtractTransactionCtx(ctx)

	if err != nil {
		logger.Error().Str("action", "extract transaction").Err(err)
		return UseCase{}, err
	}

	stockDB := &pg.StockPgxStore{DB: tx}
	logDB := &pg.PgxProvisionLog{DB: tx}
	inventoryDB := &pg.InventoryPgxStore{DB: tx}

	uc := UseCase{
		ctx:         ctx,
		stockDB:     stockDB,
		logDB:       logDB,
		inventoryDB: inventoryDB,
		logger:      &logger,
	}

	return uc, nil
}
