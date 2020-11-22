package recipe

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/middleware"
	persistence "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/persistence/postgres"
	domain "github.com/anatollupacescu/retail-sample/domain/retail/recipe"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/hlog"
)

var ErrBadItemID = errors.New("could not parse ID")

func getByID(r *http.Request) (domain.Recipe, error) {
	hlog.FromRequest(r).Info().Str("action", "enter").Msg("get recipe by id")

	vars := mux.Vars(r)
	recipeID := vars["recipeID"]

	id, err := strconv.Atoi(recipeID)

	if err != nil {
		hlog.FromRequest(r).Error().Err(err).Str("action", "convert request ID").Msg("get recipe by id")
		return domain.Recipe{}, ErrBadItemID
	}

	tx, err := middleware.ExtractTransaction(r)
	if err != nil {
		hlog.FromRequest(r).Error().Err(err).Str("action", "extract transaction").Msg("get by recipe id")
		return domain.Recipe{}, err
	}

	store := persistence.RecipePgxStore{DB: tx}

	rcp, err := store.Get(domain.ID(id))
	if err != nil {
		hlog.FromRequest(r).Error().Err(err).Str("action", "call persistence layer").Msg("get recipe by id")
		return domain.Recipe{}, err
	}

	hlog.FromRequest(r).Info().Str("action", "success").Msg("get recipe by id")

	return rcp, nil
}

func getAll(r *http.Request) ([]domain.Recipe, error) {
	hlog.FromRequest(r).Info().Str("action", "enter").Msg("list recipes")

	tx, err := middleware.ExtractTransaction(r)
	if err != nil {
		hlog.FromRequest(r).Error().Err(err).Str("action", "extract transaction").Msg("list recipes")
		return nil, err
	}

	store := persistence.RecipePgxStore{DB: tx}

	recipes, err := store.List()
	if err != nil {
		hlog.FromRequest(r).Error().Err(err).Str("action", "call persistence layer").Msg("list recipes")
		return nil, err
	}

	hlog.FromRequest(r).Info().Str("action", "success").Msg("list recipes")

	return recipes, nil
}
