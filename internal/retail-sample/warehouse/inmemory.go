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

func (i InMemoryInboundLog) Add(k time.Time, v ProvisionEntry) {
	i[k] = v
}

func (i InMemoryInboundLog) List() (r []ProvisionEntry) {
	for _, v := range i {
		r = append(r, v)
	}
	return
}

type InMemoryOutboundLog map[time.Time]SoldItem

func (m InMemoryOutboundLog) Add(i SoldItem) {
	m[i.Date] = i
}

func (m InMemoryOutboundLog) List() (r []SoldItem) {
	for _, v := range m {
		r = append(r, v)
	}
	return
}
