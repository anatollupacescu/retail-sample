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

		DB db
	}

	db interface {
		Add(string) (int, error)
		Find(string) (int, error)
		Save(*Item) error
	}

	Collection struct {
		DB db
	}
)

var (
	ErrItemNotFound  = errors.New("item not found")
	ErrEmptyName     = errors.New("name not provided")
	ErrDuplicateName = errors.New("item type already present")
)

func (i *Item) Enable() error {
	i.Enabled = true

	return i.DB.Save(i)
}

func (i *Item) Disable() error {
	i.Enabled = false

	return i.DB.Save(i)
}

func (i Collection) Add(name string) (int, error) {
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
