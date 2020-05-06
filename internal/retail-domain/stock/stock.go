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

		Inventory    Inventory
		ProvisionLog ProvisionLog
	}

	StockPosition struct {
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
		Add(ProvisionEntry) error
	}
)

var ErrItemNotFound = errors.New("stock item not found")

func (s Stock) CurrentStock() (ps []StockPosition, err error) {
	items, err := s.Inventory.List()

	if err != nil {
		return nil, err
	}

	for _, item := range items {
		itemID := int(item.ID)
		qty, err := s.Store.Quantity(itemID)

		if err == ErrItemNotFound {
			continue
		}

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

var ErrNotEnoughStock = errors.New("not enough stock")

func (s Stock) Quantity(id int) (qty int, err error) {
	return s.Store.Quantity(id)
}

func (s Stock) Provision(in []ProvisionEntry) (updatedQtys map[int]int, err error) {
	for _, i := range in {
		itemID := i.ID

		if _, err = s.Inventory.Get(itemID); err != nil {
			return nil, err
		}
	}

	updatedQtys = make(map[int]int)

	for _, spe := range in {
		newQty, err := s.Store.Provision(spe.ID, spe.Qty)

		if err != nil {
			return nil, err
		}

		updatedQtys[spe.ID] = newQty
	}

	for id, qty := range updatedQtys {
		entry := ProvisionEntry{
			ID:  id,
			Qty: qty,
		}

		if err = s.ProvisionLog.Add(entry); err != nil {
			return
		}
	}

	return updatedQtys, nil
}

func (s Stock) GetProvisionLog() ([]ProvisionEntry, error) {
	return s.ProvisionLog.List()
}

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
