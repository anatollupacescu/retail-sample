package inventory

import (
	"net/http"

	"github.com/rs/zerolog/hlog"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/middleware"
	pg "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/persistence/postgres"
	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/usecase"
)

func newUseCase(r *http.Request) (usecase.Inventory, error) {
	logger := hlog.FromRequest(r)

	tx, err := middleware.ExtractTransaction(r)

	if err != nil {
		logger.Error().Str("action", "extract transaction").Err(err)
		return usecase.Inventory{}, err
	}

	inv := tx.Inventory()
	ctx := r.Context()

	inventoryDB := &pg.InventoryPgxStore{DB: tx.Tx}

	uc := usecase.NewInventory(ctx, inv, inventoryDB)

	return uc, nil
}
