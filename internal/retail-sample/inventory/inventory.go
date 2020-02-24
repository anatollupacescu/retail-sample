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
		add(Name) ID
		find(Name) ID
		get(ID) Item
		all() []Item
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

	if i.Store.find(name) != zeroID {
		return zeroID, ErrDuplicateName
	}

	newID := i.Store.add(name)

	return newID, nil
}

func (i Inventory) All() (r []Item) {
	return i.Store.all()
}

func (i Inventory) Find(name Name) ID {
	return i.Store.find(name)
}

func (i Inventory) Get(id ID) Item {
	return i.Store.get(id)
}
