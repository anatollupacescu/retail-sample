package recipe

import domain "github.com/anatollupacescu/retail-sample/domain/retail/recipe"

func (o *UseCase) UpdateStatus(recipeID int, enabled bool) (dto domain.DTO, err error) {
	dto, err = o.recipeDB.Get(recipeID)

	if err != nil {
		return
	}

	recipe := domain.Recipe{
		ID:          dto.ID,
		Name:        dto.Name,
		Ingredients: dto.Ingredients,
		Enabled:     dto.Enabled,
		DB:          o.recipeDB,
	}

	switch enabled {
	case true:
		err = recipe.Enable()
	default:
		err = recipe.Disable()
	}

	if err != nil {
		o.logger.Error().Err(err).Msg("call domain layer")
		return
	}

	dto.Enabled = recipe.Enabled

	return
}
