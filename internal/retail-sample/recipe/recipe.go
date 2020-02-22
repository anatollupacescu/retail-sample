package recipe

import (
	"errors"

	"github.com/anatollupacescu/retail-sample/internal/retail-sample/inventory"
)

type (
	Name string
	ID   int

	Inventory interface {
		Get(int) inventory.Item
	}

	Recipe struct {
		Name        string
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
	ErrEmptyName           = errors.New("empty")
	ErrNoIngredients       = errors.New("no components found")
	ErrIgredientNotFound   = errors.New("ingredient not found")
	ErrQuantityNotProvided = errors.New("quantity not provided")
)

func (b Book) Add(name string, ingredients []Ingredient) error {
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

		if b.Inventory.Get(v.ID) == zeroItem {
			return ErrIgredientNotFound
		}
	}

	_, err := b.Store.add(Recipe{
		Name:        name,
		Ingredients: ingredients,
	})

	return err
}

func (b Book) Get(id int) Recipe {
	return b.Store.get(ID(id))
}

func (b Book) Names() (r []string) {
	for _, rp := range b.Store.all() {
		r = append(r, rp.Name)
	}

	return
}
