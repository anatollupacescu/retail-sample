package stock

import (
	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
	"github.com/anatollupacescu/retail-sample/domain/retail/recipe"

	"github.com/pkg/errors"
)

type (
	DB interface {
		Quantity(int) (int, error)
		Provision(int, int) (int, error)
		Sell([]recipe.Ingredient, int) error
	}

	Stock struct {
		DB           DB
		InventoryDB  Inventory
		ProvisionLog ProvisionLog
	}

	Inventory interface {
		List() ([]inventory.Item, error)
		Get(int) (inventory.Item, error)
	}

	ProvisionEntry struct {
		ID  int
		Qty int
	}

	ProvisionLog interface {
		Add(int, int) (int, error)
	}
)

var (
	ErrItemNotFound   = errors.New("stock item not found")
	ErrNotEnoughStock = errors.New("not enough stock")
)

func (s Stock) Provision(itemID, qty int) (id int, err error) {
	if _, err = s.InventoryDB.Get(itemID); err != nil {
		return
	}

	_, err = s.DB.Provision(itemID, qty)

	if err != nil {
		return
	}

	if id, err = s.ProvisionLog.Add(itemID, qty); err != nil {
		return
	}

	return
}

func (s Stock) Sell(ingredients []recipe.Ingredient, qty int) error {
	for _, i := range ingredients {
		presentQty, err := s.DB.Quantity(i.ID)

		switch err {
		case nil, ErrItemNotFound: //continue
		default:
			return err
		}

		if presentQty < qty*i.Qty {
			return ErrNotEnoughStock
		}
	}

	if err := s.DB.Sell(ingredients, qty); err != nil {
		return err
	}

	return nil
}
