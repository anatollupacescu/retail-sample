package stock

import (
	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
	"github.com/anatollupacescu/retail-sample/domain/retail/recipe"

	"github.com/pkg/errors"
)

type (
	db interface {
		Quantity(int) (int, error)
		Provision(int, int) error
		Sell([]recipe.Ingredient, int) error
	}

	inventoryDB interface {
		Get(int) (inventory.Item, error)
	}

	Stock struct {
		DB          db
		InventoryDB inventoryDB
	}
)

var (
	ErrItemNotFound   = errors.New("stock item not found")
	ErrNotEnoughStock = errors.New("not enough stock")
)

func (s Stock) Provision(itemID, qty int) error {
	_, err := s.InventoryDB.Get(itemID) // check it exists

	if err != nil {
		return err
	}

	return s.DB.Provision(itemID, qty)
}

func (s Stock) Sell(ingredients []recipe.Ingredient, qty int) error {
	for _, i := range ingredients {
		available, err := s.DB.Quantity(i.ID)

		switch err {
		case nil, ErrItemNotFound: //continue
		default:
			return err
		}

		if available < qty*i.Qty {
			return ErrNotEnoughStock
		}
	}

	if err := s.DB.Sell(ingredients, qty); err != nil {
		return err
	}

	return nil
}
