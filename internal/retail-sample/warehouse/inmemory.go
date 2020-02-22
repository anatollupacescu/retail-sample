package warehouse

import (
	"time"

	"github.com/anatollupacescu/retail-sample/internal/retail-sample/inventory"
	"github.com/anatollupacescu/retail-sample/internal/retail-sample/recipe"
)

func NewInMemoryStock() Stock {
	inMemoryStore := inventory.NewInMemoryStore()
	inMemoryRecipeStore := recipe.NewInMemoryStore()

	return Stock{
		inboundLog:  make(InMemoryInboundLog),
		outboundLog: make(InMemoryOutboundLog),
		recipeBook:  recipe.Book{Store: &inMemoryRecipeStore},
		inventory:   inventory.NewInventory(inMemoryStore),
	}
}

type InMemoryInboundLog map[time.Time]ProvisionEntry

func (i InMemoryInboundLog) Add(v ProvisionEntry) {
	i[time.Now()] = v
}

func (i InMemoryInboundLog) List() (r []ProvisionEntry) {
	for t, v := range i {
		r = append(r, ProvisionEntry{
			Time: t,
			ID:   v.ID,
			Qty:  v.Qty,
		})
	}
	return
}

type InMemoryOutboundLog map[time.Time]OrderLogEntry

func (m InMemoryOutboundLog) Add(i OrderLogEntry) {
	m[i.Date] = i
}

func (m InMemoryOutboundLog) List() (r []OrderLogEntry) {
	for _, v := range m {
		r = append(r, v)
	}
	return
}
