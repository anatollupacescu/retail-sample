package order

import (
	"errors"
	"time"

	"github.com/anatollupacescu/retail-sample/internal/retail-domain/recipe"
)

type (
	OrderEntry struct {
		RecipeID int
		Qty      int
	}

	ID int

	Order struct {
		ID ID
		OrderEntry
		Date time.Time
	}

	Store interface {
		Add(Order) (ID, error)
		Get(ID) (Order, error)
		List() ([]Order, error)
	}

	recipeBook interface {
		Get(recipe.ID) (recipe.Recipe, error)
	}

	orderStock interface {
		Sell(ingredients []recipe.Ingredient, qty int) error
	}

	Orders struct {
		Store      Store
		RecipeBook recipeBook
		Stock      orderStock
	}
)

var (
	ErrOrderNotFound   = errors.New("order not found")
	ErrInvalidQuantity = errors.New("quantity not valid")
	ErrInvalidRecipe   = errors.New("invalid recipe")
)

func (o Orders) PlaceOrder(id int, qty int) (orderID ID, err error) {
	var zeroOrderID ID

	if qty <= 0 {
		return zeroOrderID, ErrInvalidQuantity
	}

	recipeID := recipe.ID(id)

	r, err := o.RecipeBook.Get(recipeID)

	if err != nil {
		return zeroOrderID, err
	}

	if !r.Enabled {
		return zeroOrderID, ErrInvalidRecipe
	}

	ingredients := r.Ingredients

	if err := o.Stock.Sell(ingredients, qty); err != nil {
		return zeroOrderID, err
	}

	ord := Order{
		OrderEntry: OrderEntry{
			RecipeID: id,
			Qty:      qty,
		},
		Date: time.Now(),
	}

	orderID, err = o.Store.Add(ord)

	if err != nil {
		return zeroOrderID, err
	}

	return orderID, nil
}

func (o Orders) Get(id ID) (Order, error) {
	return o.Store.Get(id)
}

func (o Orders) List() ([]Order, error) {
	return o.Store.List()
}
