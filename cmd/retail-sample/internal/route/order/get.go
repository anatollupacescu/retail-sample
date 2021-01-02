package order

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/hlog"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/middleware"
	persistence "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/persistence/postgres"
	domain "github.com/anatollupacescu/retail-sample/domain/retail/order"
)

var ErrBadItemID = errors.New("could not parse ID")

func getByID(r *http.Request) (domain.DTO, error) {
	hlog.FromRequest(r).Info().Str("action", "enter").Msg("get order by id")

	vars := mux.Vars(r)
	orderID := vars["orderID"]

	id, err := strconv.Atoi(orderID)

	if err != nil {
		hlog.FromRequest(r).Error().Err(err).Str("action", "convert request ID").Msg("get order by id")
		return domain.DTO{}, ErrBadItemID
	}

	tx, err := middleware.ExtractTransaction(r)
	if err != nil {
		hlog.FromRequest(r).Error().Err(err).Str("action", "extract transaction").Msg("get by order id")
		return domain.DTO{}, err
	}

	store := persistence.OrderPgxStore{DB: tx}

	order, err := store.Get(id)
	if err != nil {
		hlog.FromRequest(r).Error().Err(err).Str("action", "call persistence layer").Msg("get by order id")
		return domain.DTO{}, err
	}

	hlog.FromRequest(r).Info().Str("action", "success").Msg("get by order id")

	return order, nil
}

func getAll(r *http.Request) ([]domain.DTO, error) {
	hlog.FromRequest(r).Info().Str("action", "enter").Msg("list orders")

	tx, err := middleware.ExtractTransaction(r)
	if err != nil {
		hlog.FromRequest(r).Error().Err(err).Str("action", "extract transaction").Msg("list orders")
		return nil, err
	}

	store := persistence.OrderPgxStore{DB: tx}

	orders, err := store.List()
	if err != nil {
		hlog.FromRequest(r).Error().Err(err).Str("action", "call persistence layer").Msg("list orders")
		return nil, err
	}

	hlog.FromRequest(r).Info().Str("action", "success").Msg("list orders")

	return orders, nil
}
