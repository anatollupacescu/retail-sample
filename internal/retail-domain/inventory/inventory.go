package inventory

import "errors"

type ( //inventory

	Name string
	ID   int

	Item struct {
		ID   ID
		Name Name
	}

	Store interface {
		Add(Name) ID
		Find(Name) ID
		Get(ID) Item
		List() []Item
	}

	Inventory struct {
		Store Store
	}
)

var (
	ErrEmptyName     = errors.New("name not provided")
	ErrDuplicateName = errors.New("item type already present")

	zeroID = ID(0)
)

func (i Inventory) Add(name Name) (ID, error) {
	if name == "" {
		return zeroID, ErrEmptyName
	}

	if i.Store.Find(name) != zeroID {
		return zeroID, ErrDuplicateName
	}

	newID := i.Store.Add(name)

	return newID, nil
}

func (i Inventory) List() (r []Item) {
	return i.Store.List()
}

func (i Inventory) Find(name Name) ID {
	return i.Store.Find(name)
}

func (i Inventory) Get(id ID) Item {
	return i.Store.Get(id)
}
