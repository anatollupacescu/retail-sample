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
		ID   ID
	}

	store interface {
		add(Order) ID
		all() []Order
	}

	Orders struct {
		Store store
	}
)

func (o Orders) Add(oe OrderEntry) ID {
	ord := Order{
		OrderEntry: oe,
		Date:       time.Now(),
	}

	return o.Store.add(ord)
}

func (o Orders) All() (os []Order) {
	return o.Store.all()
}
