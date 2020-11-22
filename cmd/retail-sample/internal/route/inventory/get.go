package inventory

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/hlog"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/middleware"
	persistence "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/persistence/postgres"
	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
)

func GetByID(r *http.Request) (inventory.Item, error) {
	hlog.FromRequest(r).Info().Str("action", "enter").Msg("get by id")

	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["itemID"])
	if err != nil {
		hlog.FromRequest(r).Error().Err(err).Str("action", "convert request ID").Msg("get by id")
		return inventory.Item{}, ErrParseItemID
	}

	tx, err := middleware.ExtractTransaction(r)
	if err != nil {
		hlog.FromRequest(r).Error().Err(err).Str("action", "extract transaction").Msg("get by id")
		return inventory.Item{}, err
	}

	store := persistence.InventoryPgxStore{DB: tx}

	item, err := store.Get(id)
	if err != nil {
		hlog.FromRequest(r).Error().Err(err).Str("action", "call persistence layer").Msg("get by id")
		return inventory.Item{}, err
	}

	hlog.FromRequest(r).Info().Str("action", "success").Msg("get by id")

	return item, nil
}

func ListItems(r *http.Request) ([]inventory.Item, error) {
	hlog.FromRequest(r).Info().Str("action", "enter").Msg("get all")

	tx, err := middleware.ExtractTransaction(r)
	if err != nil {
		hlog.FromRequest(r).Error().Err(err).Str("action", "extract transaction").Msg("get by id")
		return nil, err
	}

	store := persistence.InventoryPgxStore{DB: tx}

	items, err := store.List()
	if err != nil {
		hlog.FromRequest(r).Error().Err(err).Str("action", "call persistence layer").Msg("get by id")
		return nil, err
	}

	hlog.FromRequest(r).Info().Str("action", "success").Msg("get all")

	return items, nil
}
