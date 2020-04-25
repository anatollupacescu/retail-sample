package order

import "time"

type (
	OrderEntry struct {
		RecipeID int
		Qty      int
	}

	ID int

	Order struct {
		ID ID
		OrderEntry
		Date time.Time
	}

	store interface {
		Add(Order) (ID, error)
		List() ([]Order, error)
	}

	Orders struct {
		Store store
	}
)

func (o Orders) Add(oe OrderEntry) (ID, error) {
	ord := Order{
		OrderEntry: oe,
		Date:       time.Now(),
	}

	return o.Store.Add(ord)
}

func (o Orders) List() ([]Order, error) {
	return o.Store.List()
}
