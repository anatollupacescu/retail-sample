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

	Position struct {
		ID   int
		Name string
		Qty  int
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

var ErrItemNotFound = errors.New("stock item not found")

func New(store Store, inventory Inventory, log ProvisionLog) Stock {
	return Stock{
		Store:        store,
		Inventory:    inventory,
		ProvisionLog: log,
	}
}

var ErrNotEnoughStock = errors.New("not enough stock")

func (s Stock) Position(id int) (p Position, err error) {
	qty, err := s.Store.Quantity(id)

	if err != nil {
		return Position{}, err
	}

	item, err := s.Inventory.Get(id)

	if err != nil {
		return Position{}, err
	}

	p = Position{
		ID:   id,
		Name: item.Name,
		Qty:  qty,
	}

	return p, nil
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
