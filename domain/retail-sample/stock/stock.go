package stock

import (
	"github.com/anatollupacescu/retail-sample/domain/retail-sample/inventory"
	"github.com/anatollupacescu/retail-sample/domain/retail-sample/recipe"

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

func (s Stock) CurrentStock() (ps []Position, err error) {
	items, err := s.Inventory.List()

	if err != nil {
		return nil, err
	}

	for _, item := range items {
		itemID := item.ID
		qty, err := s.Store.Quantity(itemID)

		if errors.Is(err, ErrItemNotFound) {
			continue
		}

		if err != nil {
			return nil, err
		}

		ps = append(ps, Position{
			ID:   itemID,
			Name: item.Name,
			Qty:  qty,
		})
	}

	return
}

var ErrNotEnoughStock = errors.New("not enough stock")

func (s Stock) Quantity(id int) (qty int, err error) {
	return s.Store.Quantity(id)
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

func (s Stock) GetAllProvisions() ([]ProvisionEntry, error) {
	return s.ProvisionLog.List()
}

func (s Stock) GetProvision(id int) (e ProvisionEntry, err error) {
	return s.ProvisionLog.Get(id)
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
