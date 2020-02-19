package inventory

import "errors"

type ( //inventory

	Name string
	ID   int

	Record struct {
		ID   ID
		Name Name
	}

	Store interface {
		add(Name) ID
		find(Name) ID
		all() []Record
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

func (i Inventory) All() (r []Record) {
	return i.store.all()
}

func (i Inventory) Get(s string) int {
	return int(i.store.find(Name(s)))
}
