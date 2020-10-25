package recipe

import (
	"github.com/anatollupacescu/retail-sample/domain/retail/recipe"
)

func (o *Recipe) GetAll() (recipes []recipe.Recipe, err error) {
	o.logger.Info("get all", "enter")

	recipes, err = o.book.List()

	if err != nil {
		o.logger.Error("get all", "call domain layer", err)

		return
	}

	o.logger.Info("get all", "success")

	return
}
