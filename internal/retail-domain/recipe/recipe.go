package recipe

import (
	"errors"

	"github.com/anatollupacescu/retail-sample/internal/retail-domain/inventory"
)

type (
	Inventory interface {
		Get(inventory.ID) (inventory.Item, error)
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
		List() ([]Recipe, error)
		Get(ID) (Recipe, error)
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

func (b Book) Add(name Name, ingredients []Ingredient) (ID, error) {
	var zeroRecipeID = ID(0)

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
		itemID := inventory.ID(v.ID)

		_, err := b.Inventory.Get(itemID)

		switch err {
		case nil:
			continue
		case inventory.ErrInventoryItemNotFound:
			return zeroRecipeID, ErrIgredientNotFound
		default:
			return zeroRecipeID, err
		}
	}

	return b.Store.Add(Recipe{
		Name:        name,
		Ingredients: ingredients,
	})
}

var ErrRecipeNotFound = errors.New("recipe not found")

func (b Book) Get(id ID) (Recipe, error) {
	return b.Store.Get(id)
}

func (b Book) List() (r []Recipe, err error) {
	list, err := b.Store.List()

	r = append(r, list...)

	return
}
