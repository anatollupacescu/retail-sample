package recipe

import (
	"errors"

	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
)

type (
	Inventory interface {
		Get(int) (inventory.Item, error)
	}

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

	Store interface {
		Add(Recipe) (ID, error)
		List() ([]Recipe, error)
		Get(ID) (Recipe, error)
		Save(Recipe) error
	}

	Book struct {
		Store     Store
		Inventory Inventory
	}
)

var (
	ErrEmptyName           = errors.New("empty name")
	ErrNoIngredients       = errors.New("no ingredients provided")
	ErrIgredientNotFound   = errors.New("ingredient not found")
	ErrIgredientNotEnabled = errors.New("ingredient not enabled")
	ErrQuantityNotProvided = errors.New("quantity not provided")
)

func New(store Store, inventory Inventory) Book {
	return Book{
		Store:     store,
		Inventory: inventory,
	}
}

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
			return zeroRecipeID, ErrIgredientNotEnabled
		}
	}

	return b.Store.Add(Recipe{
		Name:        name,
		Ingredients: ingredients,
		Enabled:     true,
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

func (b Book) SetStatus(id ID, enabled bool) error {
	r, err := b.Store.Get(id)

	if err != nil {
		return err
	}

	r.Enabled = enabled

	err = b.Store.Save(r)

	return err
}
