package inventory

import (
	"errors"
	"strings"
)

type (
	Item struct {
		ID      int
		Name    string
		Enabled bool
	}

	Store interface {
		Add(string) (int, error)
		Find(string) (int, error)
		Get(int) (Item, error)
		List() ([]Item, error)
		Update(Item) error
	}

	Inventory struct {
		store Store
	}
)

var (
	ErrItemNotFound  = errors.New("item not found")
	ErrEmptyName     = errors.New("name not provided")
	ErrDuplicateName = errors.New("item type already present")
)

func New(store Store) Inventory {
	return Inventory{store: store}
}

func (i Inventory) UpdateStatus(id int, enabled bool) (item Item, err error) {
	item, err = i.store.Get(id)

	switch err {
	case nil: // continue
	case ErrItemNotFound:
		return
	default:
		return Item{}, err
	}

	item.Enabled = enabled

	err = i.store.Update(item)

	if err != nil {
		return Item{}, err
	}

	return item, nil
}

func (i Inventory) Add(name string) (int, error) {
	if strings.TrimSpace(name) == "" {
		return 0, ErrEmptyName
	}

	_, err := i.store.Find(name)

	switch err {
	case ErrItemNotFound: //success
		return i.store.Add(name)
	case nil:
		return 0, ErrDuplicateName
	default:
		return 0, err
	}
}
