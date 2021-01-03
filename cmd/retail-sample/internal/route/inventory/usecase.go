package inventory

import (
	"net/http"

	"github.com/rs/zerolog/hlog"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/machine/inventory"
	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/middleware"
)

func newUseCase(r *http.Request) (inventory.UseCase, error) {
	logger := hlog.FromRequest(r)

	tx, err := middleware.ExtractTransaction(r)

	if err != nil {
		logger.Error().Str("action", "extract transaction").Err(err)
		return inventory.UseCase{}, err
	}

	ctx := r.Context()

	uc := inventory.New(ctx, tx)

	return uc, nil
}
