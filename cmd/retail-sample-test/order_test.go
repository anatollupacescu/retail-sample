package acceptance_test

import (
	"errors"

	"github.com/anatollupacescu/arbortest/runner"

	domain "github.com/anatollupacescu/retail-sample/domain/retail-sample/order"

	http "github.com/anatollupacescu/retail-sample/cmd/retail-sample-test"
	random "github.com/anatollupacescu/retail-sample/cmd/retail-sample-test"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/app/order"
)

func createOrder(recipeID, qty int) (o domain.Order, err error) {
	cl := http.Post("order")

	return order.Create(recipeID, qty, cl)
}

// group:order after:recipe
func testCreateOrderOK(t *runner.T) {
	createdName := random.Name()
	item, _ := createItem(createdName)

	id := item.ID

	_, _ = provisionStock(id, 30)

	ingredients := map[int]int{
		id: 6,
	}

	name := random.Name()

	r, _ := createRecipe(name, ingredients)

	o, err := createOrder(int(r.ID), 5)

	if err != nil {
		t.Error(err)
		return
	}

	if o.ID == 0 {
		t.Error(errors.New("should have a non zero value ID"))
	}
}

// group:order
func testCreateOrderWhenNotEnoughStock(t *runner.T) {
	createdName := random.Name()
	item, _ := createItem(createdName)

	id := item.ID

	ingredients := map[int]int{
		id: 1,
	}

	name := random.Name()

	r, _ := createRecipe(name, ingredients)

	_, err := createOrder(int(r.ID), 1)

	if err == nil {
		t.Error(errors.New("should have errored"))
	}
}
