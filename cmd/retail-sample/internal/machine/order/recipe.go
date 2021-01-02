package order

import (
	pg "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/persistence/postgres"
	domain "github.com/anatollupacescu/retail-sample/domain/retail/order"
)

type validator struct {
	recipes *pg.RecipePgxStore
}

func (v *validator) Valid(id int) error {
	recipe, err := v.recipes.Get(id)

	if err != nil {
		return err
	}

	if !recipe.Enabled {
		return domain.ErrInvalidRecipe
	}

	return nil
}
