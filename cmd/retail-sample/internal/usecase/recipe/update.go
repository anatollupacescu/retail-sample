package recipe

import (
	"errors"

	"github.com/anatollupacescu/retail-sample/domain/retail/recipe"
)

var ErrBadItemID = errors.New("could not parse ID")

type UpdateStatusDTO struct {
	RecipeID int
	Enabled  bool
}

func (o *Recipe) Update(in UpdateStatusDTO) (recipe.Recipe, error) {
	o.logger.Info("update status", "begin")

	recipeID := recipe.ID(in.RecipeID)

	err := o.book.SetStatus(recipeID, in.Enabled)

	if err != nil {
		o.logger.Error("update status", "call domain", err)
		return recipe.Recipe{}, err
	}

	rec, err := o.book.Get(recipeID)

	if err != nil {
		o.logger.Error("update status", "call domain to retrieve updated entity", err)
		return recipe.Recipe{}, err
	}

	o.logger.Info("update status", "success")

	return rec, nil
}
