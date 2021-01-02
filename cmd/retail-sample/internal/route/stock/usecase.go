package stock

import (
	"net/http"

	"github.com/rs/zerolog/hlog"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/machine/stock"
	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/middleware"
)

func useCase(r *http.Request) (stock.UseCase, error) {
	logger := hlog.FromRequest(r)

	tx, err := middleware.ExtractTransaction(r)

	if err != nil {
		logger.Error().Str("action", "extract transaction").Err(err)
		return stock.UseCase{}, err
	}

	ctx := r.Context()
	uc := stock.New(ctx, tx)

	return uc, nil
}
