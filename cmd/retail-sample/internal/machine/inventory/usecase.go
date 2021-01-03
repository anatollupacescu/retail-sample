package inventory

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	pg "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/persistence/postgres"
	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
)

type UseCase struct {
	ctx         context.Context
	logger      *zerolog.Logger
	inventory   *inventory.Collection
	inventoryDB *pg.InventoryPgxStore
}

func New(ctx context.Context, t pg.TX) UseCase {
	logger := log.Ctx(ctx).With().Str("layer", "use case").Logger()
	db := &pg.InventoryPgxStore{DB: t.Tx}

	return UseCase{
		ctx:         ctx,
		inventory:   &inventory.Collection{DB: db},
		inventoryDB: db,
		logger:      &logger,
	}
}
