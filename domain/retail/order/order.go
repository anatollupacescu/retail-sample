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

	recipeValidator interface {
		Valid(int) (bool, error)
	}

	stock interface {
		Extract(int, int) error
	}

	Orders struct {
		DB              db
		RecipeValidator recipeValidator
		Stock           stock
	}
)

var (
	ErrOrderNotFound   = errors.New("order not found")
	ErrInvalidQuantity = errors.New("quantity not valid")
	ErrInvalidRecipe   = errors.New("recipe not valid")
)

func (o Orders) Create(recipeID, orderCount int) (int, error) {
	if orderCount <= 0 {
		return 0, ErrInvalidQuantity
	}

	valid, err := o.RecipeValidator.Valid(recipeID)

	if err != nil {
		return 0, err
	}

	if !valid {
		return 0, ErrInvalidRecipe
	}

	if err = o.Stock.Extract(recipeID, orderCount); err != nil {
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
