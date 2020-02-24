package warehouse

import (
	"time"

	"github.com/anatollupacescu/retail-sample/internal/retail-sample/inventory"
	"github.com/anatollupacescu/retail-sample/internal/retail-sample/recipe"
)

func NewInMemoryStock() Stock {
	inventoryStore := inventory.NewInMemoryStore()
	inventory := inventory.Inventory{Store: &inventoryStore}
	recipeStore := recipe.NewInMemoryStore()
	recipeBook := recipe.Book{Store: &recipeStore, Inventory: &inventory}
	return Stock{
		InboundLog:  make(InMemoryInboundLog),
		OutboundLog: make(InMemoryOutboundLog),
		RecipeBook:  recipeBook,
		Inventory:   inventory,
		Data:        make(map[int]int),
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
