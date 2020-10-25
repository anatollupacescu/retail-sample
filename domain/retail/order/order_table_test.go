package order_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/anatollupacescu/retail-sample/domain/retail/order"
	"github.com/anatollupacescu/retail-sample/domain/retail/recipe"
)

func TestPlaceOrderT(t *testing.T) {
	tt := []struct {
		testName string
		store    *testStore
		recipes  *testRecipes
		stock    *testStock

		id  int
		qty int

		orderID order.ID
		err     error
	}{
		{
			"rejects zero quantity",
			nil, nil, nil,

			1, 0,
			order.ID(0), order.ErrInvalidQuantity,
		}, {
			"given recipe can not be retrieved it propagates error",
			nil, &testRecipes{
				get: func(id recipe.ID) (recipe.Recipe, error) {
					assert.Equal(t, recipe.ID(1), id)
					return recipe.Recipe{}, errors.New("not found")
				},
			}, nil,

			1, 1,
			order.ID(0), errors.New("not found"),
		}, {
			"given recipe is disabled it returns error",
			nil, &testRecipes{
				get: func(id recipe.ID) (recipe.Recipe, error) {
					return recipe.Recipe{Enabled: false}, nil
				},
			}, nil,

			1, 1,
			order.ID(0), order.ErrInvalidRecipe,
		}, {
			"given an enabled recipe " +
				"when ingredients can not be sold " +
				"it errors",
			nil, &testRecipes{
				get: func(id recipe.ID) (recipe.Recipe, error) {
					return recipe.Recipe{
						Enabled: true,
						Ingredients: []recipe.Ingredient{
							{
								ID: 5,
							},
						},
					}, nil
				},
			},
			&testStock{
				sell: func(i []recipe.Ingredient, qty int) error {
					assert.Equal(t, []recipe.Ingredient{
						{
							ID: 5,
						}}, i, "wrong arguments received")
					return errors.New("expected")
				},
			},

			1, 1,
			order.ID(0), errors.New("expected"),
		}, {
			"given a valid recipe it saves it to the store",
			&testStore{
				add: func(o order.Order) (order.ID, error) {
					assert.Equal(t, 1, o.Entry.Qty)
					assert.Equal(t, 45, o.Entry.RecipeID)
					return order.ID(4), nil
				},
			},
			&testRecipes{
				get: func(id recipe.ID) (recipe.Recipe, error) {
					return recipe.Recipe{
						ID:      id,
						Enabled: true,
					}, nil
				},
			},
			&testStock{
				sell: func(i []recipe.Ingredient, qty int) error {
					return nil
				},
			},

			45, 1,
			order.ID(4), nil,
		},
	}

	for _, v := range tt {
		t.Run(v.testName, func(t *testing.T) {
			orders := order.Orders{
				Store:      v.store,
				RecipeBook: v.recipes,
				Stock:      v.stock,
			}

			orderID, err := orders.PlaceOrder(v.id, v.qty)

			assert.Equal(t, v.orderID, orderID)
			assert.Equal(t, v.err, err)
		})
	}
}

type testStore struct {
	add  func(order.Order) (order.ID, error)
	get  func(order.ID) (order.Order, error)
	list func() ([]order.Order, error)
}

func (s *testStore) Add(o order.Order) (order.ID, error) {
	return s.add(o)
}
func (s *testStore) Get(id order.ID) (order.Order, error) {
	return s.get(id)
}
func (s *testStore) List() ([]order.Order, error) {
	return s.list()
}

type testRecipes struct {
	get func(recipe.ID) (recipe.Recipe, error)
}

func (r *testRecipes) Get(rID recipe.ID) (recipe.Recipe, error) {
	return r.get(rID)
}

type testStock struct {
	sell func([]recipe.Ingredient, int) error
}

func (ts *testStock) Sell(ingredients []recipe.Ingredient, qty int) error {
	return ts.sell(ingredients, qty)
}
