package order

import (
	"net/http"

	"github.com/rs/zerolog/hlog"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/machine/order"
	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/middleware"
)

func newUseCase(r *http.Request) (order.UseCase, error) {
	logger := hlog.FromRequest(r)

	tx, err := middleware.ExtractTransaction(r)

	if err != nil {
		logger.Error().Str("action", "extract transaction").Err(err)
		return order.UseCase{}, err
	}

	ctx := r.Context()

	uc := order.NewUseCase(ctx, tx)

	return uc, nil
}
