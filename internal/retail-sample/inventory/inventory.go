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
		store Store
	}
)

func NewInventory(s Store) Inventory {
	return Inventory{
		store: s,
	}
}

var (
	ErrEmptyName     = errors.New("name not provided")
	ErrDuplicateName = errors.New("item type already present")
)

var zero = 0

func (i Inventory) Add(s string) (int, error) {
	if s == "" {
		return zero, ErrEmptyName
	}

	if i.store.find(Name(s)) != ID(zero) {
		return zero, ErrDuplicateName
	}

	newID := i.store.add(Name(s))

	return int(newID), nil
}

func (i Inventory) All() (r []Item) {
	return i.store.all()
}

func (i Inventory) Find(s string) int {
	return int(i.store.find(Name(s)))
}

func (i Inventory) Get(id int) Item {
	return i.store.get(ID(id))
}
