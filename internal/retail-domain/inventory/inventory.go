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
		Store Store
	}
)

var (
	ErrItemNotFound = errors.New("item not found")

	ErrEmptyName     = errors.New("name not provided")
	ErrDuplicateName = errors.New("item type already present")
)

func (i Inventory) UpdateStatus(id int, enabled bool) (item Item, err error) {
	item, err = i.Store.Get(id)

	switch err {
	case nil:
		break
	case ErrItemNotFound:
		return
	default:
		return Item{}, err
	}

	item.Enabled = enabled

	err = i.Store.Update(item)

	return item, err
}

func (i Inventory) Add(name string) (int, error) {
	if strings.TrimSpace(name) == "" {
		return 0, ErrEmptyName
	}

	_, err := i.Store.Find(name)

	switch err {
	case ErrItemNotFound: //success
		break
	case nil:
		return 0, ErrDuplicateName
	default:
		return 0, err
	}

	return i.Store.Add(name)
}

func (i Inventory) List() ([]Item, error) {
	return i.Store.List()
}

func (i Inventory) Find(name string) (int, error) {
	return i.Store.Find(name)
}

func (i Inventory) Get(id int) (Item, error) {
	return i.Store.Get(id)
}
