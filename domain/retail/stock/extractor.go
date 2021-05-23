package stock

import (
	"github.com/pkg/errors"

	"github.com/anatollupacescu/retail-sample/domain/retail/recipe"
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
	if count <= 0 {
		return ErrInvalidExtractQuantity
	}

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

		if err := pos.extract(totalQty); err != nil {
			return err
		}
	}

	return nil
}
