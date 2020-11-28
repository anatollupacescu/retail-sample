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
	}

	Ingredient struct {
		ID  int
		Qty int
	}

	db interface {
		Add(Recipe) (ID, error)
		List() ([]Recipe, error)
		Get(ID) (Recipe, error)
		Save(Recipe) error
	}

	Inventory interface {
		Get(int) (inventory.Item, error)
	}

	Book struct {
		DB        db
		Inventory Inventory
	}
)

var (
	ErrEmptyName           = errors.New("empty name")
	ErrNoIngredients       = errors.New("no ingredients provided")
	ErrIgredientNotFound   = errors.New("ingredient not found")
	ErrIgredientDisabled   = errors.New("ingredient disabled")
	ErrQuantityNotProvided = errors.New("quantity not provided")
)

func (b Book) Add(name Name, ingredients []Ingredient) (ID, error) {
	var zeroRecipeID ID

	if name == "" {
		return zeroRecipeID, ErrEmptyName
	}

	if len(ingredients) == 0 {
		return zeroRecipeID, ErrNoIngredients
	}

	for _, v := range ingredients {
		if v.Qty == 0 {
			return zeroRecipeID, ErrQuantityNotProvided
		}
	}

	for _, v := range ingredients {
		itemID := v.ID

		item, err := b.Inventory.Get(itemID)

		switch err {
		case nil: //continue
		case inventory.ErrItemNotFound:
			return zeroRecipeID, ErrIgredientNotFound
		default:
			return zeroRecipeID, err
		}

		if !item.Enabled {
			return zeroRecipeID, ErrIgredientDisabled
		}
	}

	return b.DB.Add(Recipe{
		Name:        name,
		Ingredients: ingredients,
		Enabled:     true,
	})
}

var ErrRecipeNotFound = errors.New("recipe not found")

func (b Book) UpdateStatus(id ID, enabled bool) error {
	r, err := b.DB.Get(id)

	if err != nil {
		return err
	}

	r.Enabled = enabled

	err = b.DB.Save(r)

	return err
}
