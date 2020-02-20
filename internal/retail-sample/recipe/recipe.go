package recipe

import (
	"errors"
)

type (
	Name string
	ID   int

	Inventory interface {
		Get(int) string
	}

	Book struct {
		Store     Store
		Inventory Inventory
	}

	Recipe struct {
		Name string
	}

	Store interface {
		add(Recipe) (ID, error)
		all() []Recipe
		get(ID) Recipe
	}

	Ingredient struct {
		ID  int
		Qty int
	}
)

var (
	ErrEmptyName           = errors.New("empty")
	ErrNoIngredients       = errors.New("no components found")
	ErrIgredientNotFound   = errors.New("ingredient not found")
	ErrQuantityNotProvided = errors.New("quantity not provided")
)

var (
	zero = ""
)

func (b Book) Add(name string, ingredients []Ingredient) error {
	if name == "" {
		return ErrEmptyName
	}

	if len(ingredients) == 0 {
		return ErrNoIngredients
	}

	for _, v := range ingredients {
		if v.Qty == 0 {
			return ErrQuantityNotProvided
		}

		if b.Inventory.Get(v.ID) == zero {
			return ErrIgredientNotFound
		}
	}

	_, err := b.Store.add(Recipe{
		Name: name,
	})

	return err
}
