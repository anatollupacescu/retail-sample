package stock

import (
	"net/http"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/usecase"

	"github.com/rs/zerolog/hlog"

	pg "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/persistence/postgres"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/middleware"
)

func useCase(r *http.Request) (usecase.Stock, error) {
	logger := hlog.FromRequest(r)

	tx, err := middleware.ExtractTransaction(r)

	if err != nil {
		logger.Error().Str("action", "extract transaction").Err(err)
		return usecase.Stock{}, err
	}

	stock := tx.Stock()

	stockDB := &pg.StockPgxStore{DB: tx.Tx}
	inventoryDB := &pg.InventoryPgxStore{DB: tx.Tx}
	logDB := &pg.PgxProvisionLog{DB: tx.Tx}

	ctx := r.Context()
	uc := usecase.NewStock(ctx, stock, stockDB, logDB, inventoryDB)

	return uc, nil
}
