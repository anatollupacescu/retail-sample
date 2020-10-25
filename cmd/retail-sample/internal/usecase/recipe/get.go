package recipe

import (
	"errors"
	"strconv"

	"github.com/anatollupacescu/retail-sample/domain/retail/recipe"
)

var ErrBadItemID = errors.New("could not parse ID")

func (o *Recipe) GetByID(recipeID string) (r recipe.Recipe, err error) {
	o.logger.Info("get by id", "enter")

	var id int

	id, err = strconv.Atoi(recipeID)

	if err != nil {
		o.logger.Error("get by id", "convert request ID", err)
		return recipe.Recipe{}, ErrBadItemID
	}

	r, err = o.book.Get(recipe.ID(id))

	if err != nil {
		o.logger.Error("get by id", "call domain layer", err)

		return
	}

	o.logger.Info("get by id", "success")

	return r, nil
}
