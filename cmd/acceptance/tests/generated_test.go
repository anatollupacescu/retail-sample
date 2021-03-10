package tests

import (
	"testing"

	arbor "github.com/anatollupacescu/arbortest/runner"
)

func TestArbor(t *testing.T) {
	g := arbor.New()

	t.Run("inventory", func(t *testing.T) {
		at := arbor.NewT(t)
		g.Group("inventory")
		g.Append(at, "CreateInventoryItem", testCreateInventoryItem)
		g.Append(at, "DisableItem", testDisableItem)
		g.Append(at, "EnableItem", testEnableItem)
	})

	t.Run("recipe", func(t *testing.T) {
		at := arbor.NewT(t)
		g.Group("recipe")
		g.After(at, "inventory")
		g.Append(at, "CreateRecipe", testCreateRecipe)
		g.Append(at, "DisableRecipe", testDisableRecipe)
		g.Append(at, "EnableRecipe", testEnableRecipe)
	})

	t.Run("stock", func(t *testing.T) {
		at := arbor.NewT(t)
		g.Group("stock")
		g.After(at, "inventory")
		g.Append(at, "Provision", testProvision)
		g.Append(at, "Extract", testExtract)
	})

	t.Run("order", func(t *testing.T) {
		at := arbor.NewT(t)
		g.Group("order")
		g.After(at, "recipe")
		g.Append(at, "CreateOrder", testCreateOrder)
	})

	output := g.JSON()

	arbor.Upload(output)
}
