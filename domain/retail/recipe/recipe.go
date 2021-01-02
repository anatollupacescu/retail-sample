package recipe

import (
	"errors"
)

type (
	DTO struct {
		ID          int
		Name        string
		Ingredients []InventoryItem
		Enabled     bool
	}

	Recipe struct {
		ID          int
		Name        string
		Ingredients []InventoryItem
		Enabled     bool

		DB db
	}

	InventoryItem struct {
		ID  int
		Qty int
	}

	db interface {
		Add(DTO) (int, error)
		Find(string) (*DTO, error)
		Save(*DTO) error
	}

	inventory interface {
		Validate(...int) error
	}

	Recipes struct {
		DB        db
		Inventory inventory
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

func checkPreconditions(name string, ingredients []InventoryItem) error {
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

func (c Recipes) Add(name string, ingredients []InventoryItem) (int, error) {
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

	var ids = make([]int, 0, len(ingredients))
	for _, i := range ingredients {
		ids = append(ids, i.ID)
	}

	err = c.Inventory.Validate(ids...)
	if err != nil {
		return 0, err
	}

	dto := DTO{
		Name:        name,
		Ingredients: ingredients,
		Enabled:     true,
	}

	return c.DB.Add(dto)
}

var ErrRecipeNotFound = errors.New("recipe not found")

func (r *Recipe) Disable() error {
	dto := DTO{
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
	dto := DTO{
		ID: r.ID, Name: r.Name, Ingredients: r.Ingredients, Enabled: true,
	}

	err := r.DB.Save(&dto)

	if err != nil {
		return err
	}

	r.Enabled = true

	return nil
}
