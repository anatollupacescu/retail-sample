package order

import (
	"errors"
	"net/http"

	usecase "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/usecase/order"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/middleware"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

var ErrCreateFail = errors.New("create use case failed")

func useCase(r *http.Request) (usecase.Order, error) {
	logger := hlog.FromRequest(r)

	tx, err := middleware.ExtractTransaction(r)

	if err != nil {
		logger.Error().Str("action", "extract transaction").Err(err)
		return usecase.Order{}, err
	}

	logWrapper := logWrapper{
		logger: logger,
	}

	orders := tx.Orders()
	ctx := r.Context()

	uc := usecase.New(ctx, orders, logWrapper)

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
