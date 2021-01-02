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

	OrderDTO struct {
		ID ID
		Entry
		Date time.Time
	}

	db interface {
		Add(OrderDTO) (ID, error)
	}

	recipes interface {
		Get(recipe.ID) (recipe.RecipeDTO, error)
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

func (o Orders) Add(recipeID int, orderCount int) (orderID ID, err error) {
	if orderCount <= 0 {
		return 0, ErrInvalidQuantity
	}

	id := recipe.ID(recipeID)

	recipe, err := o.Recipes.Get(id)

	if err != nil {
		return 0, err
	}

	if !recipe.Enabled {
		return 0, ErrInvalidRecipe
	}

	for _, ingredient := range recipe.Ingredients {
		inventoryID := ingredient.ID
		totalQty := ingredient.Qty * orderCount

		err := o.Stock.Extract(inventoryID, totalQty)

		if err != nil {
			return 0, err
		}
	}

	ord := OrderDTO{
		Entry: Entry{
			RecipeID: recipeID,
			Qty:      orderCount,
		},
		Date: time.Now(),
	}

	orderID, err = o.DB.Add(ord)

	if err != nil {
		return 0, err
	}

	return orderID, nil
}
