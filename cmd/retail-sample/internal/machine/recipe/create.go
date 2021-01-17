package recipe

import (
	"github.com/pkg/errors"

	usecase "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/machine"
	domain "github.com/anatollupacescu/retail-sample/domain/retail/recipe"
)

func (o *UseCase) Create(name string, items []domain.InventoryItem) (domain.DTO, error) {
	var err error
	defer func() {
		if err != nil {
			o.logger.Error().Str("action", "create order").Err(err).Send()
		}
	}()

	id, err := o.recipes.Create(name, items)

	switch {
	case err == nil:
	case errors.Is(err, domain.ErrIngredientNotFound):
		return domain.DTO{}, errors.Wrap(usecase.ErrNotFound, err.Error())
	case
		err == domain.ErrEmptyName,
		err == domain.ErrQuantityNotProvided,
		err == domain.ErrNoIngredients:
		return domain.DTO{}, errors.Wrap(usecase.ErrBadRequest, err.Error())
	}

	o.logger.Info().Int("id", id).Msg("successfully created recipe")

	recipe, err := o.recipeDB.Get(id)

	if err != nil {
		return domain.DTO{}, err
	}

	return recipe, nil
}
