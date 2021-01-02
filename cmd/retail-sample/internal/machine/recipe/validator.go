package recipe

import (
	pg "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/persistence/postgres"
	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
	"github.com/anatollupacescu/retail-sample/domain/retail/recipe"
)

type validator struct {
	Inventory *pg.InventoryPgxStore
}

func (a validator) Validate(ids ...int) error {
	for _, id := range ids {
		item, err := a.Inventory.Get(id)

		switch err {
		case nil: //continue
		case inventory.ErrItemNotFound:
			return recipe.ErrIgredientNotFound
		default:
			return err
		}

		if !item.Enabled {
			return recipe.ErrIgredientDisabled
		}
	}

	return nil
}
