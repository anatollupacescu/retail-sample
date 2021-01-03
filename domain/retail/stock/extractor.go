package stock

import (
	"github.com/anatollupacescu/retail-sample/domain/retail/recipe"
)

type recipes interface {
	Get(int) (recipe.DTO, error)
}

type Extractor struct {
	Stock   db
	Recipes recipes
}

func (o Extractor) Extract(recipeID int, count int) error {
	recipe, err := o.Recipes.Get(recipeID)
	if err != nil {
		return err
	}

	for _, ingredient := range recipe.Ingredients {
		inventoryID := ingredient.ID
		totalQty := ingredient.Qty * count

		dto, err := o.Stock.Get(inventoryID)
		if err != nil {
			return err
		}

		pos := Position{
			InventoryID: dto.InventoryID,
			Qty:         dto.Qty,
			DB:          o.Stock,
		}

		if err := pos.Extract(totalQty); err != nil {
			return err
		}
	}

	return nil
}
