package stock

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/hlog"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/machine/stock"
	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/middleware"
	persistence "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/persistence/postgres"
)

var ErrBadItemID = errors.New("could not parse ID")

func getByID(r *http.Request) (stock.Position, error) {
	hlog.FromRequest(r).Info().Str("action", "enter").Msg("get stock position by id")

	vars := mux.Vars(r)
	itemID := vars["itemID"]

	id, err := strconv.Atoi(itemID)
	if err != nil {
		hlog.FromRequest(r).Error().Err(err).Str("action", "convert request ID").Msg("get stock position by id")
		return stock.Position{}, ErrBadItemID
	}

	tx, err := middleware.ExtractTransaction(r)
	if err != nil {
		hlog.FromRequest(r).Error().Err(err).Str("action", "extract transaction").Msg("get by recipe id")
		return stock.Position{}, err
	}

	store := persistence.InventoryPgxStore{DB: tx}

	item, err := store.Get(id)
	if err != nil {
		hlog.FromRequest(r).Error().Err(err).Str("action", "extract transaction").Msg("get by recipe id")
		return stock.Position{}, err
	}

	stockDB := persistence.StockPgxStore{DB: tx}

	dto, err := stockDB.Get(id)

	if err != nil {
		hlog.FromRequest(r).Error().Err(err).Str("action", "call persistence layer").Msg("get stock position by id")
		return stock.Position{}, err
	}

	hlog.FromRequest(r).Info().Str("action", "success").Msg("get stock position by id")

	pos := stock.Position{
		ID:   id,
		Name: item.Name,
		Qty:  dto.Qty,
	}

	return pos, nil
}

func getAll(r *http.Request) ([]stock.Position, error) {
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

	stockDB := persistence.StockPgxStore{DB: tx}

	positions := make([]stock.Position, 0, len(items))

	for i := range items {
		item := items[i]
		dto, err := stockDB.Get(item.ID)

		switch err {
		case nil, persistence.ErrStockItemNotFound:
		default:
			hlog.FromRequest(r).Error().Err(err).Str("action", "find item quantity").Msg("list stock positions")
		}

		pos := stock.Position{
			ID:   item.ID,
			Name: item.Name,
			Qty:  dto.Qty,
		}

		positions = append(positions, pos)
	}

	hlog.FromRequest(r).Info().Str("action", "success").Msg("list recipes")

	return positions, nil
}

func getProvisionLog(r *http.Request) ([]persistence.ProvisionEntry, error) {
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
