package recipe

import (
	"errors"

	"github.com/anatollupacescu/retail-sample/internal/retail-sample/inventory"
)

type (
	Inventory interface {
		Get(inventory.ID) inventory.Item
	}

	Name string
	ID   int

	Recipe struct {
		Name        Name
		Ingredients []Ingredient
	}

	Store interface {
		add(Recipe) (ID, error)
		all() []Recipe
		get(ID) Recipe
	}

	Book struct {
		Store     Store
		Inventory Inventory
	}

	Ingredient struct {
		ID  int
		Qty int
	}
)

var (
	ErrEmptyName           = errors.New("empty name")
	ErrNoIngredients       = errors.New("no components found")
	ErrIgredientNotFound   = errors.New("ingredient not found")
	ErrQuantityNotProvided = errors.New("quantity not provided")
)

func (b Book) Add(name Name, ingredients []Ingredient) error {
	if name == "" {
		return ErrEmptyName
	}

	if len(ingredients) == 0 {
		return ErrNoIngredients
	}

	var zeroItem inventory.Item

	for _, v := range ingredients {
		if v.Qty == 0 {
			return ErrQuantityNotProvided
		}

		itemID := inventory.ID(v.ID)

		if b.Inventory.Get(itemID) == zeroItem {
			return ErrIgredientNotFound
		}
	}

	_, err := b.Store.add(Recipe{
		Name:        name,
		Ingredients: ingredients,
	})

	return err
}

func (b Book) Get(id ID) Recipe {
	return b.Store.get(id)
}

func (b Book) Names() (r []Name) {
	for _, rp := range b.Store.all() {
		r = append(r, rp.Name)
	}

	return
}
