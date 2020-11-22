package inventory

import (
	"errors"
	"net/http"

	usecase "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/usecase/inventory"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/middleware"
	pg "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/persistence/postgres"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

var ErrCreateFail = errors.New("create use case failed")

func useCase(r *http.Request) (usecase.Inventory, error) {
	logger := hlog.FromRequest(r)

	tx, err := middleware.ExtractTransaction(r)

	if err != nil {
		logger.Error().Str("action", "extract transaction").Err(err)
		return usecase.Inventory{}, err
	}

	logWrapper := logWrapper{
		logger: logger,
	}

	inv := tx.Inventory()
	ctx := r.Context()

	inventoryDB := &pg.InventoryPgxStore{DB: tx.Tx}

	uc := usecase.New(ctx, inv, inventoryDB, logWrapper)

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
