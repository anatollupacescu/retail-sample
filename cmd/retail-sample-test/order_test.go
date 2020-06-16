package acceptance_test

import (
	"errors"

	domain "github.com/anatollupacescu/retail-sample/internal/retail-domain/order"

	http "github.com/anatollupacescu/retail-sample/cmd/retail-sample-test"
	random "github.com/anatollupacescu/retail-sample/cmd/retail-sample-test"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/app/order"
)

func createOrder(recipeID, qty int) (o domain.Order, err error) {
	cl := http.Post("order")

	return order.Create(recipeID, qty, cl)
}

func testCreateOrderOK() error {
	id := createRandomItem()

	_, _ = provisionStock(id, 30)

	ingredients := map[int]int{
		id: 6,
	}

	name := random.Name()

	r, _ := createRecipe(name, ingredients)

	o, err := createOrder(int(r.ID), 5)

	if err != nil {
		return err
	}

	if o.ID == 0 {
		return errors.New("should have a non zero value ID")
	}

	return nil
}

func testCreateOrderWhenNotEnoughStock() error {
	id := createRandomItem()

	ingredients := map[int]int{
		id: 1,
	}

	name := random.Name()

	r, _ := createRecipe(name, ingredients)

	_, err := createOrder(int(r.ID), 1)

	if err == nil {
		return errors.New("should have rejected")
	}

	return nil
}
