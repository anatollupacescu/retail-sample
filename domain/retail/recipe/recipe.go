package recipe

import (
	"errors"

	inv "github.com/anatollupacescu/retail-sample/domain/retail/inventory"
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
		Get(int) (DTO, error)
		Add(DTO) (int, error)
		Find(string) (DTO, error)
		Save(DTO) error
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

func (c Recipes) Create(name string, ingredients []InventoryItem) (int, error) {
	if err := checkPreconditions(name, ingredients); err != nil {
		return 0, err
	}

	err := checkIsAlreadyPresent(c.DB, name)
	if err != nil {
		return 0, err
	}

	err = checkIngredientsAreValid(c.Inventory, ingredients)
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

func (r *Recipe) Disable() error {
	dto := DTO{
		ID:          r.ID,
		Name:        r.Name,
		Ingredients: r.Ingredients,

		Enabled: false,
	}

	err := r.DB.Save(dto)

	if err != nil {
		return err
	}

	r.Enabled = false

	return nil
}

func (r *Recipe) Enable() error {
	dto := DTO{
		ID:          r.ID,
		Name:        r.Name,
		Ingredients: r.Ingredients,

		Enabled: true,
	}

	err := r.DB.Save(dto)

	if err != nil {
		return err
	}

	r.Enabled = true

	return nil
}

var (
	ErrDuplicateName  = errors.New("duplicate name")
	ErrRecipeNotFound = errors.New("recipe not found")
)

func checkIsAlreadyPresent(db db, name string) error {
	_, err := db.Find(name)
	switch err {
	case ErrRecipeNotFound:
		return nil
	case nil:
		return ErrDuplicateName
	default:
		return err
	}
}

var ErrIngredientNotFound = errors.New("ingredient not found")

func checkIngredientsAreValid(validator inventory, ingredients []InventoryItem) error {
	var ids = make([]int, 0, len(ingredients))
	for _, i := range ingredients {
		ids = append(ids, i.ID)
	}

	err := validator.Validate(ids...)

	switch err {
	case nil:
		return nil
	case inv.ErrItemNotFound:
		return ErrIngredientNotFound
	default:
		return err
	}
}
