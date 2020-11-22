package stock

import (
	"errors"
	"net/http"

	usecase "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/usecase/stock"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/middleware"
	postgres "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/persistence/postgres"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

var ErrCreateFail = errors.New("create use case failed")

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

	stock := tx.Stock()
	proLog := tx.ProvisionLog()

	inventoryDB := &postgres.InventoryPgxStore{DB: tx.Tx}
	stockDB := &postgres.StockPgxStore{DB: tx.Tx}

	uc := usecase.New(r.Context(), stock, proLog, stockDB, inventoryDB, logWrapper)

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
