package recipe

import (
	"strconv"

	"github.com/anatollupacescu/retail-sample/domain/retail/recipe"
)

func (o *Recipe) Update(reqID string, enabled bool) (r recipe.Recipe, err error) {
	o.logger.Info("get all", "enter")

	var id int

	id, err = strconv.Atoi(reqID)

	if err != nil {
		o.logger.Error("get by id", "convert request ID", err)
		return recipe.Recipe{}, ErrBadItemID
	}

	recipeID := recipe.ID(id)
	err = o.book.SetStatus(recipeID, enabled)

	if err != nil {
		o.logger.Error("get all", "call domain layer", err)
		return
	}

	r, err = o.book.Get(recipeID)

	if err != nil {
		o.logger.Error("get all", "call domain layer to retrieve the newly updated recipe", err)
		return
	}

	o.logger.Info("get all", "success")

	return
}
