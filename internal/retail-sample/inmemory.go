package retailsample

import (
	"errors"

	"github.com/anatollupacescu/retail-sample/internal/retail-domain/recipe"

	"time"
)

//provision log
type InMemoryProvisionLog map[time.Time]ProvisionEntry

func (i InMemoryProvisionLog) Add(v ProvisionEntry) {
	i[time.Now()] = v
}

func (i InMemoryProvisionLog) List() (r []ProvisionEntry) {
	for _, v := range i {
		r = append(r, ProvisionEntry{
			ID:  v.ID,
			Qty: v.Qty,
		})
	}
	return
}

// stock
type InMemoryStock struct {
	data map[int]int
}

func NewStockWithData(data map[int]int) Stock {
	return InMemoryStock{
		data: data,
	}
}

func NewInMemoryStock() Stock {
	return InMemoryStock{
		data: make(map[int]int),
	}
}

func (s InMemoryStock) Quantity(id int) int {
	return s.data[id]
}

func (s InMemoryStock) Provision(id, qty int) int {
	newQty := s.data[id] + qty

	s.data[id] = newQty

	return newQty
}

var ErrNotEnoughStock = errors.New("not enough stock")

func (s InMemoryStock) Sell(ii []recipe.Ingredient, qty int) error {
	for _, i := range ii {
		presentQty := s.data[i.ID]
		requestedQty := i.Qty * qty
		if requestedQty > presentQty {
			return ErrNotEnoughStock
		}
	}

	for _, i := range ii {
		presentQty := s.data[i.ID]
		requestedQty := i.Qty * qty
		s.data[i.ID] = presentQty - requestedQty
	}

	return nil
}
