package order

import (
	"errors"
	"time"
)

type (
	DTO struct {
		ID       int
		RecipeID int
		Qty      int
		Date     time.Time
	}

	db interface {
		Add(DTO) (int, error)
	}

	recipes interface {
		Valid(int) error
	}

	stock interface {
		Extract(int, int) error
	}

	Orders struct {
		DB      db
		Recipes recipes
		Stock   stock
	}
)

var (
	ErrOrderNotFound   = errors.New("order not found")
	ErrInvalidQuantity = errors.New("quantity not valid")
	ErrInvalidRecipe   = errors.New("invalid recipe")
)

func (o Orders) Add(recipeID, orderCount int) (int, error) {
	if orderCount <= 0 {
		return 0, ErrInvalidQuantity
	}

	if err := o.Recipes.Valid(recipeID); err != nil {
		return 0, err
	}

	if err := o.Stock.Extract(recipeID, orderCount); err != nil {
		return 0, err
	}

	ord := DTO{
		RecipeID: recipeID,
		Qty:      orderCount,
		Date:     time.Now(),
	}

	orderID, err := o.DB.Add(ord)

	if err != nil {
		return 0, err
	}

	return orderID, nil
}
