package recipe

import (
	"github.com/anatollupacescu/retail-sample/domain/retail/recipe"
	domain "github.com/anatollupacescu/retail-sample/domain/retail/recipe"
)

type CreateDTO struct {
	Name        domain.Name
	Ingredients []domain.Ingredient
}

func (o *Recipe) Create(dto CreateDTO) (recipe recipe.Recipe, err error) {
	o.logger.Info("get all", "enter")

	id, err := o.book.Add(dto.Name, dto.Ingredients)

	if err != nil {
		o.logger.Error("get all", "call domain layer", err)
		return
	}

	recipe, err = o.book.Get(id)

	if err != nil {
		o.logger.Error("get all", "call domain layer to retrieve the newly created recipe", err)
		return
	}

	o.logger.Info("get all", "success")

	return
}
