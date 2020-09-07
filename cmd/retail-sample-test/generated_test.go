// +build acceptance

package acceptance_test

import (
	"testing"

	arbor "github.com/anatollupacescu/arbortest/runner"
)

func TestArbor(t *testing.T) {
	g := arbor.New()

	t.Run("inventory", func(t *testing.T) {
		at := arbor.NewT(t)
		g.Group("inventory")
		g.Append(at, "Create", testCreate)
		g.Append(at, "CreateWithEmptyName", testCreateWithEmptyName)
		g.Append(at, "Duplicate", testDuplicate)
		g.Append(at, "Disable", testDisable)
		g.Append(at, "GetAll", testGetAll)
		g.Append(at, "GetOne", testGetOne)
	})

	t.Run("recipe", func(t *testing.T) {
		at := arbor.NewT(t)
		g.Group("recipe")
		g.After(at, "inventory")
		g.Append(at, "CreateRecipe", testCreateRecipe)
		g.Append(at, "CreateRecipeNoName", testCreateRecipeNoName)
		g.Append(at, "CreateRecipeNoItems", testCreateRecipeNoItems)
		g.Append(at, "GetRecipe", testGetRecipe)
		g.Append(at, "GetAllRecipes", testGetAllRecipes)
		g.Append(at, "DisableRecipe", testDisableRecipe)
	})

	t.Run("stock", func(t *testing.T) {
		at := arbor.NewT(t)
		g.Group("stock")
		g.After(at, "inventory")
		g.Append(at, "Provision", testProvision)
		g.Append(at, "GetStockPos", testGetStockPos)
		g.Append(at, "GetAllStockPos", testGetAllStockPos)
	})

	t.Run("order", func(t *testing.T) {
		at := arbor.NewT(t)
		g.Group("order")
		g.After(at, "recipe")
		g.Append(at, "CreateOrderOK", testCreateOrderOK)
		g.Append(at, "CreateOrderWhenNotEnoughStock", testCreateOrderWhenNotEnoughStock)
	})

	output := g.JSON()

	arbor.Upload(output)
}
