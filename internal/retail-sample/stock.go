package retailsample

import (
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/inventory"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/recipe"
	"github.com/pkg/errors"
)

type (
	stock struct {
		store Stock
	}

	StockPosition struct {
		ID   int
		Name string
		Qty  int
	}
)

func (s *stock) CurrentStock(ii []inventory.Item) (ps []StockPosition) {
	for _, item := range ii {
		itemID := int(item.ID)
		qty, err := s.store.Quantity(itemID)

		if err != nil {
			return nil
		}

		ps = append(ps, StockPosition{
			ID:   itemID,
			Name: string(item.Name),
			Qty:  qty,
		})
	}

	return
}

func (s *stock) Sell(ingredients []recipe.Ingredient, qty int) error {
	for _, i := range ingredients {
		presentQty, err := s.store.Quantity(i.ID)

		if err != nil {
			return err
		}

		if presentQty < qty*i.Qty {
			return errors.Wrapf(BusinessErr, ErrNotEnoughStock.Error())
		}
	}

	if err := s.store.Sell(ingredients, qty); err != nil {
		return err
	}

	return nil
}
