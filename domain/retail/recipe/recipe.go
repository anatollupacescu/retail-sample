package recipe

import (
	"github.com/pkg/errors"

	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
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

	inventoryValidator interface {
		Valid(int) (bool, error)
	}

	Recipes struct {
		DB            db
		ItemValidator inventoryValidator
	}
)

func (c Recipes) Create(name string, ingredients []InventoryItem) (int, error) {
	if err := checkPreconditions(name, ingredients); err != nil {
		return 0, err
	}

	if err := checkIsAlreadyPresent(c.DB, name); err != nil {
		return 0, err
	}

	if err := checkIngredientsAreValid(c.ItemValidator, ingredients); err != nil {
		return 0, err
	}

	dto := DTO{
		Name:        name,
		Ingredients: ingredients,
		Enabled:     true,
	}

	return c.DB.Add(dto)
}

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
			return errors.Wrapf(ErrQuantityNotProvided, "id: %d", v.ID)
		}
	}

	return nil
}

var (
	ErrDuplicateName = errors.New("duplicate name")
	ErrNotFound      = errors.New("recipe not found")
)

func checkIsAlreadyPresent(db db, name string) error {
	_, err := db.Find(name)

	switch err {
	case ErrNotFound:
		return nil
	case nil:
		return ErrDuplicateName
	default:
		return err
	}
}

var (
	ErrIngredientNotFound = errors.New("ingredient not found")
	ErrIngredientNotValid = errors.New("ingredient not valid")
)

func checkIngredientsAreValid(validator inventoryValidator, ingredients []InventoryItem) error {
	for _, i := range ingredients {
		valid, err := validator.Valid(i.ID)

		if err == nil {
			if !valid {
				return errors.Wrapf(ErrIngredientNotValid, "id: %d", i.ID)
			}

			continue
		}

		if errors.Is(err, inventory.ErrNotFound) {
			return errors.Wrapf(ErrIngredientNotFound, "id not found: %d", i.ID)
		}

		return err
	}

	return nil
}

func (r *Recipe) Disable() error {
	dto := DTO{
		ID:          r.ID,
		Name:        r.Name,
		Ingredients: r.Ingredients,

		Enabled: false,
	}

	if err := r.DB.Save(dto); err != nil {
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

	if err := r.DB.Save(dto); err != nil {
		return err
	}

	r.Enabled = true

	return nil
}
