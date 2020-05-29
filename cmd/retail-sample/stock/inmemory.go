package stock

import (
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/recipe"
	domain "github.com/anatollupacescu/retail-sample/internal/retail-domain/stock"

	"time"
)

type InMemoryProvisionLog map[time.Time]domain.ProvisionEntry

func (i InMemoryProvisionLog) Add(v domain.ProvisionEntry) error {
	i[time.Now()] = v

	return nil
}

func (i InMemoryProvisionLog) List() (r []domain.ProvisionEntry, err error) {
	for _, v := range i {
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
