package order

import (
	pg "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/persistence/postgres"
	"github.com/anatollupacescu/retail-sample/domain/retail/stock"
)

type extractor struct {
	stock   *pg.StockPgxStore
	recipes *pg.RecipePgxStore
}

func (o extractor) Extract(recipeID int, count int) error {
	recipe, err := o.recipes.Get(recipeID)
	if err != nil {
		return err
	}

	for _, ingredient := range recipe.Ingredients {
		inventoryID := ingredient.ID
		totalQty := ingredient.Qty * count

		dto, err := o.stock.Get(inventoryID)
		if err != nil {
			return err
		}

		pos := stock.Position{
			InventoryID: dto.InventoryID,
			Qty:         dto.Qty,
			DB:          o.stock,
		}

		if err := pos.Extract(totalQty); err != nil {
			return err
		}
	}

	return nil
}
