package stock

import (
	"net/http"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/usecase"

	"github.com/rs/zerolog"
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

	logWrapper := logWrapper{
		logger: logger,
	}

	ctx := r.Context()

	stock := tx.Stock()

	stockDB := &pg.StockPgxStore{DB: tx.Tx}
	inventoryDB := &pg.InventoryPgxStore{DB: tx.Tx}

	uc := usecase.NewStock(ctx, stock, stockDB, inventoryDB, logWrapper)

	return uc, nil
}

type logWrapper struct {
	logger *zerolog.Logger
}

func (l logWrapper) Error(action, message string, err error) {
	l.logger.Error().Str("action", action).Err(err).Msg(message)
}

func (l logWrapper) Info(action, message string) {
	l.logger.Info().Str("action", action).Msg(message)
}
