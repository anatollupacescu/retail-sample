package stock

import (
	"github.com/anatollupacescu/retail-sample/domain/retail/recipe"
	"github.com/pkg/errors"
)

type recipes interface {
	Get(int) (recipe.DTO, error)
}

type Extractor struct {
	Stock   db
	Recipes recipes
}

var ErrRecipeDisabled = errors.New("can not process disabled recipe")

func (e Extractor) Extract(recipeID int, count int) error {
	recipe, err := e.Recipes.Get(recipeID)
	if err != nil {
		return err
	}

	if !recipe.Enabled {
		return ErrRecipeDisabled
	}

	for _, ingredient := range recipe.Ingredients {
		inventoryID := ingredient.ID
		totalQty := ingredient.Qty * count

		dto, err := e.Stock.Get(inventoryID)

		switch err {
		case nil:
		case ErrPositionNotFound:
			return ErrNotEnoughStock
		default:
			return err
		}

		pos := Position{
			InventoryID: dto.InventoryID,
			Qty:         dto.Qty,
			DB:          e.Stock,
		}

		if err := pos.Extract(totalQty); err != nil {
			return err
		}
	}

	return nil
}
