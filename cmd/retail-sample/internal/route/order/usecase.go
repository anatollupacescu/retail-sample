package order

import (
	"net/http"

	"github.com/rs/zerolog/hlog"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/middleware"
	pg "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/persistence/postgres"
	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/usecase"
)

func newUseCase(r *http.Request) (usecase.Order, error) {
	logger := hlog.FromRequest(r)

	tx, err := middleware.ExtractTransaction(r)

	if err != nil {
		logger.Error().Str("action", "extract transaction").Err(err)
		return usecase.Order{}, err
	}

	orders := tx.Orders()
	ctx := r.Context()

	orderDB := &pg.OrderPgxStore{DB: tx.Tx}

	uc := usecase.NewOrder(ctx, orders, orderDB)

	return uc, nil
}
