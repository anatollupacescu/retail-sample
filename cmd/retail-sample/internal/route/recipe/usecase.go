package recipe

import (
	"net/http"

	"github.com/rs/zerolog/hlog"

	usecase "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/machine/recipe"
	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/middleware"
)

func newUseCase(r *http.Request) (usecase.UseCase, error) {
	logger := hlog.FromRequest(r)

	tx, err := middleware.ExtractTransaction(r)

	if err != nil {
		logger.Error().Str("action", "extract transaction").Err(err)
		return usecase.UseCase{}, err
	}

	ctx := r.Context()

	uc := usecase.New(ctx, tx)

	return uc, nil
}
