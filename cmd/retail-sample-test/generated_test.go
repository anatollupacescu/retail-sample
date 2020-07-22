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
		g.Append(at, "testCreate", testCreate)
		g.Append(at, "testCreateWithEmptyName", testCreateWithEmptyName)
		g.Append(at, "testDuplicate", testDuplicate)
		g.Append(at, "testDisable", testDisable)
		g.Append(at, "testGetAll", testGetAll)
		g.Append(at, "testGetOne", testGetOne)
	})

	t.Run("recipe", func(t *testing.T) {
		at := arbor.NewT(t)
		g.Group("recipe")
		g.After(at, "inventory")
		g.Append(at, "testCreateRecipe", testCreateRecipe)
		g.Append(at, "testCreateRecipeNoName", testCreateRecipeNoName)
		g.Append(at, "testCreateRecipeNoItems", testCreateRecipeNoItems)
		g.Append(at, "testGetRecipe", testGetRecipe)
		g.Append(at, "testGetAllRecipes", testGetAllRecipes)
		g.Append(at, "testDisableRecipe", testDisableRecipe)
	})

	t.Run("stock", func(t *testing.T) {
		at := arbor.NewT(t)
		g.Group("stock")
		g.After(at, "inventory")
		g.Append(at, "testProvision", testProvision)
		g.Append(at, "testGetStockPos", testGetStockPos)
		g.Append(at, "testGetAllStockPos", testGetAllStockPos)
	})

	t.Run("order", func(t *testing.T) {
		at := arbor.NewT(t)
		g.Group("order")
		g.After(at, "recipe")
		g.Append(at, "testCreateOrderOK", testCreateOrderOK)
		g.Append(at, "testCreateOrderWhenNotEnoughStock", testCreateOrderWhenNotEnoughStock)
	})

	output := g.JSON()

	arbor.Upload(output)
}
