package inventory

import (
	"errors"
	"strings"
)

type (
	ItemDTO struct {
		ID      int
		Name    string
		Enabled bool
	}

	Item struct {
		ID      int
		Name    string
		Enabled bool

		DB db
	}

	db interface {
		Add(string) (int, error)
		Find(string) (int, error)
		Save(*ItemDTO) error
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
	dto := ItemDTO{
		ID: i.ID, Name: i.Name, Enabled: true,
	}

	if err := i.DB.Save(&dto); err != nil {
		return err
	}

	i.Enabled = true

	return nil
}

func (i *Item) Disable() error {
	dto := ItemDTO{
		ID: i.ID, Name: i.Name, Enabled: false,
	}

	if err := i.DB.Save(&dto); err != nil {
		return err
	}

	i.Enabled = false

	return nil
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
