package stock

import (
	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
	"github.com/anatollupacescu/retail-sample/domain/retail/recipe"

	"github.com/pkg/errors"
)

type (
	Store interface {
		Quantity(int) (int, error)
		Provision(int, int) (int, error)
		Sell([]recipe.Ingredient, int) error
	}

	Stock struct {
		Store Store

		Inventory    Inventory
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
		List() ([]ProvisionEntry, error)
		Add(int, int) (int, error)
		Get(int) (ProvisionEntry, error)
	}
)

var (
	ErrItemNotFound   = errors.New("stock item not found")
	ErrNotEnoughStock = errors.New("not enough stock")
)

func New(store Store, inventory Inventory, log ProvisionLog) Stock {
	return Stock{
		Store:        store,
		Inventory:    inventory,
		ProvisionLog: log,
	}
}

func (s Stock) Provision(itemID, qty int) (id int, err error) {
	if _, err = s.Inventory.Get(itemID); err != nil {
		return
	}

	_, err = s.Store.Provision(itemID, qty)

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
		presentQty, err := s.Store.Quantity(i.ID)

		switch err {
		case nil, ErrItemNotFound: //continue
		default:
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
