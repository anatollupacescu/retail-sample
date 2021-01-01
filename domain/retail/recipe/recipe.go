package recipe

import (
	"errors"

	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
)

type (
	Name string
	ID   int

	RecipeDTO struct {
		ID          ID
		Name        Name
		Ingredients []InventoryItem
		Enabled     bool
	}

	Recipe struct {
		ID          ID
		Name        Name
		Ingredients []InventoryItem
		Enabled     bool

		DB db
	}

	InventoryItem struct {
		ID  int
		Qty int
	}

	db interface {
		Add(RecipeDTO) (ID, error)
		Find(Name) (*RecipeDTO, error)
		Save(*RecipeDTO) error
	}

	Inventory interface {
		Get(int) (inventory.ItemDTO, error)
	}

	Recipes struct {
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

func checkPreconditions(name Name, ingredients []InventoryItem) error {
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

func (c Recipes) Add(name Name, ingredients []InventoryItem) (ID, error) {
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

	for _, i := range ingredients {
		item, err := c.Inventory.Get(i.ID)

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

	dto := RecipeDTO{
		Name:        name,
		Ingredients: ingredients,
		Enabled:     true,
	}

	return c.DB.Add(dto)
}

var ErrRecipeNotFound = errors.New("recipe not found")

func (r *Recipe) Disable() error {
	dto := RecipeDTO{
		ID: r.ID, Name: r.Name, Ingredients: r.Ingredients, Enabled: false,
	}

	err := r.DB.Save(&dto)

	if err != nil {
		return err
	}

	r.Enabled = false

	return nil
}

func (r *Recipe) Enable() error {
	dto := RecipeDTO{
		ID: r.ID, Name: r.Name, Ingredients: r.Ingredients, Enabled: true,
	}

	err := r.DB.Save(&dto)

	if err != nil {
		return err
	}

	r.Enabled = true

	return nil
}
