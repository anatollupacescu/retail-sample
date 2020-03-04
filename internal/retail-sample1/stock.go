package retailsampleapp1

import (
	"errors"

	"github.com/anatollupacescu/retail-sample/internal/retail-domain/recipe"
)

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
