package inmemory

import (
	"errors"

	"github.com/anatollupacescu/retail-sample/domain/retail-sample/recipe"
	"github.com/anatollupacescu/retail-sample/domain/retail-sample/stock"
)

type ProvisionLog struct {
	serial *int
	data   map[int]stock.ProvisionEntry
}

func NewProvisionLog() stock.ProvisionLog {
	return &ProvisionLog{
		serial: new(int),
		data:   make(map[int]stock.ProvisionEntry),
	}
}

func (i ProvisionLog) Add(itemID, qty int) (int, error) {
	*i.serial++

	id := *i.serial

	i.data[id] = stock.ProvisionEntry{
		ID:  itemID,
		Qty: qty,
	}

	return id, nil
}

func (i ProvisionLog) Get(id int) (e stock.ProvisionEntry, err error) {
	var ok bool

	if e, ok = i.data[id]; ok {
		return e, nil
	}

	return e, errors.New("not found")
}

func (i ProvisionLog) List() (r []stock.ProvisionEntry, err error) {
	for _, v := range i.data {
		r = append(r, stock.ProvisionEntry{
			ID:  v.ID,
			Qty: v.Qty,
		})
	}

	return
}

// Stock store.
type Stock struct {
	data map[int]int
}

func NewStock() Stock {
	return Stock{
		data: make(map[int]int),
	}
}

func (s Stock) Quantity(id int) (int, error) {
	return s.data[id], nil
}

func (s Stock) Provision(id, qty int) (int, error) {
	newQty := s.data[id] + qty

	s.data[id] = newQty

	return newQty, nil
}

func (s Stock) Sell(ii []recipe.Ingredient, qty int) error {
	for _, i := range ii {
		presentQty := s.data[i.ID]
		requestedQty := i.Qty * qty

		if requestedQty > presentQty {
			return stock.ErrNotEnoughStock
		}
	}

	for _, i := range ii {
		presentQty := s.data[i.ID]
		requestedQty := i.Qty * qty
		s.data[i.ID] = presentQty - requestedQty
	}

	return nil
}
