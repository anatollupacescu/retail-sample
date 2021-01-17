package recipe

import (
	"strconv"

	"github.com/pkg/errors"

	usecase "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/machine"
	domain "github.com/anatollupacescu/retail-sample/domain/retail/recipe"
)

func (o *UseCase) GetByID(recipeID string) (domain.DTO, error) {
	var err error
	defer func() {
		if err != nil {
			o.logger.Error().Str("action", "get order").Err(err).Send()
		}
	}()

	id, err := strconv.Atoi(recipeID)

	if err != nil {
		return domain.DTO{}, errors.Wrapf(usecase.ErrBadRequest, "parse recipe ID: %v", recipeID)
	}

	rcp, err := o.recipeDB.Get(id)

	switch err {
	case nil:
	case domain.ErrNotFound:
		return domain.DTO{}, errors.Wrapf(usecase.ErrNotFound, "get order with id: %v", id)
	default:
		return domain.DTO{}, err
	}

	return rcp, nil
}

func (o *UseCase) GetAll() ([]domain.Recipe, error) {
	recipes, err := o.recipeDB.List()

	if err != nil {
		o.logger.Error().Str("action", "get all").Err(err).Send()
		return nil, err
	}

	return recipes, nil
}
