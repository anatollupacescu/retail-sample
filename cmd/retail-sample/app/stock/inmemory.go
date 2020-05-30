package stock

import (
	"errors"

	"github.com/anatollupacescu/retail-sample/internal/retail-domain/recipe"
	domain "github.com/anatollupacescu/retail-sample/internal/retail-domain/stock"
)

type InMemoryProvisionLog struct {
	serial *int
	data   map[int]domain.ProvisionEntry
}

func NewInMemoryProvisionLog() domain.ProvisionLog {
	return &InMemoryProvisionLog{
		serial: new(int),
		data:   make(map[int]domain.ProvisionEntry),
	}
}

func (i InMemoryProvisionLog) Add(itemID, qty int) (int, error) {
	*i.serial++

	id := *i.serial

	i.data[id] = domain.ProvisionEntry{
		ID:  itemID,
		Qty: qty,
	}

	return id, nil
}

func (i InMemoryProvisionLog) Get(id int) (e domain.ProvisionEntry, err error) {
	var ok bool

	if e, ok = i.data[id]; ok {
		return e, nil
	}

	return e, errors.New("not found")
}

func (i InMemoryProvisionLog) List() (r []domain.ProvisionEntry, err error) {
	for _, v := range i.data {
		r = append(r, domain.ProvisionEntry{
			ID:  v.ID,
			Qty: v.Qty,
		})
	}

	return
}

//InMemoryStore store
type InMemoryStore struct {
	data map[int]int
}

func NewInMemoryStock() domain.StockStore {
	return InMemoryStore{
		data: make(map[int]int),
	}
}

func (s InMemoryStore) Quantity(id int) (int, error) {
	return s.data[id], nil
}

func (s InMemoryStore) Provision(id, qty int) (int, error) {
	newQty := s.data[id] + qty

	s.data[id] = newQty

	return newQty, nil
}

func (s InMemoryStore) Sell(ii []recipe.Ingredient, qty int) error {
	for _, i := range ii {
		presentQty := s.data[i.ID]
		requestedQty := i.Qty * qty
		if requestedQty > presentQty {
			return domain.ErrNotEnoughStock
		}
	}

	for _, i := range ii {
		presentQty := s.data[i.ID]
		requestedQty := i.Qty * qty
		s.data[i.ID] = presentQty - requestedQty
	}

	return nil
}
