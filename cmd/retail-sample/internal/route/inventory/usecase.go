package inventory

import (
	"net/http"

	"github.com/rs/zerolog/hlog"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/machine/inventory"
	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/middleware"
)

func newUseCase(r *http.Request) (inventory.Inventory, error) {
	logger := hlog.FromRequest(r)

	tx, err := middleware.ExtractTransaction(r)

	if err != nil {
		logger.Error().Str("action", "extract transaction").Err(err)
		return inventory.Inventory{}, err
	}

	ctx := r.Context()

	uc := inventory.New(ctx, tx)

	return uc, nil
}
