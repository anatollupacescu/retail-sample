package stock

import (
	"errors"
	"net/http"
	"strconv"

	usecase2 "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/usecase"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/hlog"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/middleware"
	persistence "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/persistence/postgres"
	domain "github.com/anatollupacescu/retail-sample/domain/retail/stock"
)

var ErrBadItemID = errors.New("could not parse ID")

func getByID(r *http.Request) (usecase2.Position, error) {
	hlog.FromRequest(r).Info().Str("action", "enter").Msg("get stock position by id")

	vars := mux.Vars(r)
	itemID := vars["itemID"]

	id, err := strconv.Atoi(itemID)
	if err != nil {
		hlog.FromRequest(r).Error().Err(err).Str("action", "convert request ID").Msg("get stock position by id")
		return usecase2.Position{}, ErrBadItemID
	}

	tx, err := middleware.ExtractTransaction(r)
	if err != nil {
		hlog.FromRequest(r).Error().Err(err).Str("action", "extract transaction").Msg("get by recipe id")
		return usecase2.Position{}, err
	}

	store := persistence.InventoryPgxStore{DB: tx}

	item, err := store.Get(id)
	if err != nil {
		hlog.FromRequest(r).Error().Err(err).Str("action", "extract transaction").Msg("get by recipe id")
		return usecase2.Position{}, err
	}

	stock := persistence.StockPgxStore{DB: tx}

	qty, err := stock.Quantity(id)

	if err != nil {
		hlog.FromRequest(r).Error().Err(err).Str("action", "call persistence layer").Msg("get stock position by id")
		return usecase2.Position{}, err
	}

	hlog.FromRequest(r).Info().Str("action", "success").Msg("get stock position by id")

	pos := usecase2.Position{
		ID:   id,
		Name: item.Name,
		Qty:  qty,
	}

	return pos, nil
}

func getAll(r *http.Request) ([]usecase2.Position, error) {
	hlog.FromRequest(r).Info().Str("action", "enter").Msg("list stock positions")

	tx, err := middleware.ExtractTransaction(r)
	if err != nil {
		hlog.FromRequest(r).Error().Err(err).Str("action", "extract transaction").Msg("list stock positions")
		return nil, err
	}

	store := persistence.InventoryPgxStore{DB: tx}

	items, err := store.List()

	if err != nil {
		hlog.FromRequest(r).Error().Err(err).Str("action", "call persistence layer").Msg("list stock positions")
		return nil, err
	}

	stock := persistence.StockPgxStore{DB: tx}

	positions := make([]usecase2.Position, 0, len(items))

	for _, item := range items {
		qty, err := stock.Quantity(item.ID)

		switch err {
		case nil, domain.ErrItemNotFound:
		default:
			hlog.FromRequest(r).Error().Err(err).Str("action", "find item quantity").Msg("list stock positions")
		}

		pos := usecase2.Position{
			ID:   item.ID,
			Name: item.Name,
			Qty:  qty,
		}
		positions = append(positions, pos)
	}

	hlog.FromRequest(r).Info().Str("action", "success").Msg("list recipes")

	return positions, nil
}

func getProvisionLog(r *http.Request) ([]domain.ProvisionEntry, error) {
	hlog.FromRequest(r).Info().Str("action", "enter").Msg("list provision log")

	tx, err := middleware.ExtractTransaction(r)
	if err != nil {
		hlog.FromRequest(r).Error().Err(err).Str("action", "extract transaction").Msg("list provision log")
		return nil, err
	}

	stock := persistence.PgxProvisionLog{DB: tx}

	items, err := stock.List()
	if err != nil {
		hlog.FromRequest(r).Error().Err(err).Str("action", "call domain").Msg("list provision log")
		return nil, err
	}

	hlog.FromRequest(r).Info().Str("action", "success").Msg("list provision log")

	return items, nil
}
