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

	DB interface {
		Add(string) (int, error)
		Find(string) (int, error)
		Get(int) (Item, error)
		Update(Item) error
	}

	Inventory struct {
		DB DB
	}
)

var (
	ErrItemNotFound  = errors.New("item not found")
	ErrEmptyName     = errors.New("name not provided")
	ErrDuplicateName = errors.New("item type already present")
)

func (i Inventory) UpdateStatus(id int, enabled bool) error {
	item, err := i.DB.Get(id)

	if err != nil {
		return err
	}

	item.Enabled = enabled

	err = i.DB.Update(item)

	if err != nil {
		return err
	}

	return nil
}

func (i Inventory) Add(name string) (int, error) {
	if strings.TrimSpace(name) == "" {
		return 0, ErrEmptyName
	}

	_, err := i.DB.Find(name)

	switch err {
	case ErrItemNotFound: //continue
	case nil:
		return 0, ErrDuplicateName
	default:
		return 0, err
	}

	id, err := i.DB.Add(name)

	if err != nil {
		return 0, err
	}

	return id, nil
}
