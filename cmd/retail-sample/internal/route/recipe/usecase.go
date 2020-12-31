package recipe

import (
	"net/http"

	"github.com/rs/zerolog/hlog"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/middleware"
	pg "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/persistence/postgres"
	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/usecase"
)

func newUseCase(r *http.Request) (usecase.Recipe, error) {
	logger := hlog.FromRequest(r)

	tx, err := middleware.ExtractTransaction(r)

	if err != nil {
		logger.Error().Str("action", "extract transaction").Err(err)
		return usecase.Recipe{}, err
	}

	recipe := tx.Recipe()
	ctx := r.Context()

	recipeDB := &pg.RecipePgxStore{DB: tx.Tx}
	uc := usecase.NewRecipe(ctx, recipe, recipeDB)

	return uc, nil
}
