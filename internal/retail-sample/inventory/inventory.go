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

	zero = 0
)

func (i Inventory) Add(s string) (int, error) {
	if s == "" {
		return zero, ErrEmptyName
	}

	if i.Store.find(Name(s)) != ID(zero) {
		return zero, ErrDuplicateName
	}

	newID := i.Store.add(Name(s))

	return int(newID), nil
}

func (i Inventory) All() (r []Item) {
	return i.Store.all()
}

func (i Inventory) Find(s string) int {
	return int(i.Store.find(Name(s)))
}

func (i Inventory) Get(id int) Item {
	return i.Store.get(ID(id))
}
