package order

import "time"

type (
	OrderEntry struct {
		RecipeID int
		Qty      int
	}

	ID int

	Order struct {
		OrderEntry
		Date time.Time
		id   ID
	}

	store interface {
		add(OrderEntry) ID
		all() []OrderEntry
	}

	Orders struct {
		Store store
	}
)

func (o Orders) Add(oe OrderEntry) ID {
	return o.Store.add(oe)
}

func (o Orders) Get(ID) Order {
	return Order{}
}

func (Orders) All() (os []Order) {
	return
}
