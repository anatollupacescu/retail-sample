package order

import (
	"errors"
	"time"

	"github.com/anatollupacescu/retail-sample/domain/retail/recipe"
)

type (
	Entry struct {
		RecipeID int
		Qty      int
	}

	ID int

	Order struct {
		ID ID
		Entry
		Date time.Time
	}

	db interface {
		Add(Order) (ID, error)
	}

	recipes interface {
		Get(recipe.ID) (recipe.Recipe, error)
	}

	orderStock interface {
		Sell(ingredients []recipe.Ingredient, qty int) error
	}

	Orders struct {
		DB      db
		Recipes recipes
		Stock   orderStock
	}
)

var (
	ErrOrderNotFound   = errors.New("order not found")
	ErrInvalidQuantity = errors.New("quantity not valid")
	ErrInvalidRecipe   = errors.New("invalid recipe")
)

func (o Orders) Add(id int, qty int) (orderID ID, err error) {
	if qty <= 0 {
		return 0, ErrInvalidQuantity
	}

	recipeID := recipe.ID(id)

	r, err := o.Recipes.Get(recipeID)

	if err != nil {
		return 0, err
	}

	if !r.Enabled {
		return 0, ErrInvalidRecipe
	}

	err = o.Stock.Sell(r.Ingredients, qty)
	if err != nil {
		return 0, err
	}

	ord := Order{
		Entry: Entry{
			RecipeID: id,
			Qty:      qty,
		},
		Date: time.Now(),
	}

	orderID, err = o.DB.Add(ord)

	if err != nil {
		return 0, err
	}

	return orderID, nil
}
