package recipe

import (
	"github.com/anatollupacescu/retail-sample/domain/retail/recipe"
)

func (o *Recipe) Create(name recipe.Name, ingredients []recipe.Ingredient) (recipe recipe.Recipe, err error) {
	o.logger.Info("get all", "enter")

	id, err := o.book.Add(name, ingredients)

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
