package recipe

import (
	"strconv"

	"github.com/pkg/errors"

	usecase "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/machine"
	"github.com/anatollupacescu/retail-sample/domain/retail/recipe"
)

func (o *UseCase) UpdateStatus(recipeID string, enabled bool) (recipe.DTO, error) {
	var err error

	defer func() {
		if err != nil {
			o.logger.Error().Str("action", "update").Err(err).Send()
		}
	}()

	id, err := strconv.Atoi(recipeID)

	if err != nil {
		return recipe.DTO{}, errors.Wrapf(usecase.ErrBadRequest, "parse recipe ID: %s", recipeID)
	}

	dto, err := o.recipeDB.Get(id)

	switch err {
	case nil:
	case recipe.ErrNotFound:
		return recipe.DTO{}, errors.Wrapf(usecase.ErrNotFound, "get recipe with id %d: %v", id, err)
	default:
		return recipe.DTO{}, err
	}

	r := recipe.Recipe{
		ID:          dto.ID,
		Name:        dto.Name,
		Ingredients: dto.Ingredients,
		Enabled:     dto.Enabled,
		DB:          o.recipeDB,
	}

	switch enabled {
	case true:
		err = r.Enable()
	default:
		err = r.Disable()
	}

	if err != nil {
		return recipe.DTO{}, err
	}

	dto.Enabled = r.Enabled

	o.logger.Info().Int("id", id).Msg("successfully updated recipe")

	return dto, nil
}
