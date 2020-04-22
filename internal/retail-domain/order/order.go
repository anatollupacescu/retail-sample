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
		Add(Order) ID
		List() []Order
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

	return o.Store.Add(ord)
}

func (o Orders) List() (os []Order) {
	return o.Store.List()
}
