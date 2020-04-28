package stock

import (
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/inventory"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/recipe"
	"github.com/pkg/errors"
)

type (
	StockStore interface {
		Quantity(int) (int, error)
		Provision(int, int) (int, error)
		Sell([]recipe.Ingredient, int) error
	}

	Stock struct {
		Store StockStore
	}

	StockPosition struct {
		ID   int
		Name string
		Qty  int
	}
)

func (s Stock) CurrentStock(ii []inventory.Item) (ps []StockPosition, err error) {
	for _, item := range ii {
		itemID := int(item.ID)
		qty, err := s.Store.Quantity(itemID)

		if err != nil {
			return nil, err
		}

		ps = append(ps, StockPosition{
			ID:   itemID,
			Name: string(item.Name),
			Qty:  qty,
		})
	}

	return
}

var ErrNotEnoughStock = errors.New("not enough stock to sell")

func (s Stock) Sell(ingredients []recipe.Ingredient, qty int) error {
	for _, i := range ingredients {
		presentQty, err := s.Store.Quantity(i.ID)

		if err != nil {
			return err
		}

		if presentQty < qty*i.Qty {
			return ErrNotEnoughStock
		}
	}

	if err := s.Store.Sell(ingredients, qty); err != nil {
		return err
	}

	return nil
}

type StockProvisionEntry struct {
	ID  int
	Qty int
}

func (s Stock) Quantity(id int) (qty int, err error) {
	return s.Store.Quantity(id)
}

func (s Stock) Provision(in []StockProvisionEntry) (map[int]int, error) {
	out := make(map[int]int, 0)

	for _, spe := range in {
		newQty, err := s.Store.Provision(spe.ID, spe.Qty)

		if err != nil {
			return nil, err
		}

		out[spe.ID] = newQty
	}

	return out, nil
}
