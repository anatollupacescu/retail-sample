package order

import (
	"errors"
	"time"

	"github.com/anatollupacescu/retail-sample/internal/retail-domain/recipe"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/stock"
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

	Orders struct {
		Store      Store
		RecipeBook recipe.Book
		Stock      stock.Stock
	}
)

var ErrOrderNotFound = errors.New("order not found")

func (o Orders) PlaceOrder(id int, qty int) (orderID ID, err error) {
	recipeID := recipe.ID(id)

	r, err := o.RecipeBook.Get(recipeID)

	var noOrderID ID

	if err != nil {
		return noOrderID, err
	}

	ingredients := r.Ingredients

	if err := o.Stock.Sell(ingredients, qty); err != nil {
		return 0, err
	}

	orderID, err = o.Add(OrderEntry{
		RecipeID: id,
		Qty:      qty,
	})

	if err != nil {
		return noOrderID, err
	}

	return orderID, nil
}

func (o Orders) Get(id ID) (Order, error) {
	return o.Store.Get(id)
}

func (o Orders) Add(oe OrderEntry) (ID, error) {
	ord := Order{
		OrderEntry: oe,
		Date:       time.Now(),
	}

	return o.Store.Add(ord)
}

func (o Orders) List() ([]Order, error) {
	return o.Store.List()
}
