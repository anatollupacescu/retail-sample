package recipe

import domain "github.com/anatollupacescu/retail-sample/domain/retail/recipe"

func (o *UseCase) Create(name string, items []domain.InventoryItem) (recipe domain.DTO, err error) {
	// call domain to create recipe
	id, err := o.recipes.Create(name, items)
	if err != nil {
		o.logger.Error().Err(err).Msg("call domain layer")
		return
	}

	recipe, err = o.recipeDB.Get(id)
	if err != nil {
		o.logger.Error().Err(err).Msg("retrieve the newly created recipe")
		return
	}

	return
}
