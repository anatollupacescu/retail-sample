package recipe

import (
	"errors"

	"github.com/anatollupacescu/retail-sample/internal/retail-domain/inventory"
)

type (
	Inventory interface {
		Get(inventory.ID) inventory.Item
	}

	Name string
	ID   int

	Recipe struct {
		ID          ID
		Name        Name
		Ingredients []Ingredient
	}

	Ingredient struct {
		ID  int
		Qty int
	}

	Store interface {
		Add(Recipe) (ID, error)
		List() []Recipe
		Get(ID) Recipe
	}

	Book struct {
		Store     Store
		Inventory Inventory
	}
)

var (
	ErrEmptyName           = errors.New("empty name")
	ErrNoIngredients       = errors.New("no components found")
	ErrIgredientNotFound   = errors.New("ingredient not found")
	ErrQuantityNotProvided = errors.New("quantity not provided")
)

var (
	zeroRecipeID = ID(0)

	zeroItem inventory.Item
)

func (b Book) Add(name Name, ingredients []Ingredient) (ID, error) {
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

		itemID := inventory.ID(v.ID)

		if b.Inventory.Get(itemID) == zeroItem {
			return zeroRecipeID, ErrIgredientNotFound
		}
	}

	return b.Store.Add(Recipe{
		Name:        name,
		Ingredients: ingredients,
	})
}

func (b Book) Get(id ID) Recipe {
	return b.Store.Get(id)
}

func (b Book) List() (r []Recipe) {
	r = append(r, b.Store.List()...)

	return
}
