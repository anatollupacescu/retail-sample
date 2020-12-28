package recipe

import (
	"errors"

	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
)

type (
	Name string
	ID   int

	Recipe struct {
		ID          ID
		Name        Name
		Ingredients []Ingredient
		Enabled     bool

		DB db
	}

	Ingredient struct {
		ID  int
		Qty int
	}

	db interface {
		Add(Recipe) (ID, error)
		Find(Name) (*Recipe, error)
		Save(*Recipe) error
	}

	Inventory interface {
		Get(int) (inventory.Item, error)
	}

	Collection struct {
		DB        db
		Inventory Inventory
	}
)

var (
	ErrEmptyName           = errors.New("empty name")
	ErrNoIngredients       = errors.New("no ingredients provided")
	ErrIgredientNotFound   = errors.New("ingredient not found")
	ErrIgredientDisabled   = errors.New("ingredient disabled")
	ErrDuplicateName       = errors.New("duplicate name")
	ErrQuantityNotProvided = errors.New("quantity not provided")
)

func checkPreconditions(name Name, ingredients []Ingredient) error {
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
	}

	return nil
}

func (c Collection) Add(name Name, ingredients []Ingredient) (ID, error) {
	if err := checkPreconditions(name, ingredients); err != nil {
		return 0, err
	}

	_, err := c.DB.Find(name)

	switch err {
	case ErrRecipeNotFound: //continue
	case nil:
		return 0, ErrDuplicateName
	default:
		return 0, err
	}

	for _, v := range ingredients {
		itemID := v.ID

		item, err := c.Inventory.Get(itemID)

		switch err {
		case nil: //continue
		case inventory.ErrItemNotFound:
			return 0, ErrIgredientNotFound
		default:
			return 0, err
		}

		if !item.Enabled {
			return 0, ErrIgredientDisabled
		}
	}

	return c.DB.Add(Recipe{
		Name:        name,
		Ingredients: ingredients,
		Enabled:     true,
	})
}

var ErrRecipeNotFound = errors.New("recipe not found")

func (r *Recipe) Disable() error {
	r.Enabled = false

	return r.DB.Save(r)
}

func (r *Recipe) Enable() error {
	r.Enabled = true

	return r.DB.Save(r)
}
