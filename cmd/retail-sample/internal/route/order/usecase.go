package order

import (
	"errors"
	"net/http"

	usecase "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/usecase/order"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/middleware"
	postgres "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/persistence/postgres"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

var ErrCreateFail = errors.New("create use case failed")

func useCase(r *http.Request) (usecase.Order, error) {
	logger := hlog.FromRequest(r)

	ctxTransaction := r.Context().Value(middleware.TxKey)

	if ctxTransaction == nil {
		logger.Err(ErrCreateFail).Msg("transaction not found")

		return usecase.Order{}, ErrCreateFail
	}

	var (
		tx postgres.TX
		ok bool
	)

	if tx, ok = ctxTransaction.(postgres.TX); !ok {
		logger.Err(ErrCreateFail).Msg("transaction of a bad type")

		return usecase.Order{}, ErrCreateFail
	}

	orders := tx.Orders()

	logWrapper := logWrapper{
		logger: logger,
	}

	uc := usecase.New(r.Context(), orders, logWrapper)

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
