// +build acceptance

package acceptance_test

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/anatollupacescu/retail-sample/internal/arbor"
)

var arborURL = flag.String("arborURL", "", "graph server URL")

func TestAcceptance(t *testing.T) {
	creteInventoryItem := arbor.New("create inventory item", testCreate)
	createEmpty := arbor.New("reject empty name", testCreateWithEmptyName, creteInventoryItem)

	getOne := arbor.New("get single item", testGetOne, creteInventoryItem)
	getAll := arbor.New("get all items", testGetAll, creteInventoryItem)
	noDuplicate := arbor.New("rejects duplicate item name", testDuplicate, creteInventoryItem)
	disable := arbor.New("disable item", testDisable, creteInventoryItem)

	//stock
	provision := arbor.New("provision stock", testProvision, creteInventoryItem)
	getOneSP := arbor.New("get single stock position", testGetStockPos, provision)
	getAllSP := arbor.New("get all stock positions", testGetAllStockPos, provision)

	//recipe
	recipeAll := arbor.Alias("create recipe ingredient", creteInventoryItem)

	createRecipeOk := arbor.New("can create recipe", testCreateRecipe, recipeAll)
	createRecipeReject := arbor.New("reject empty recipe name", testCreateRecipeNoName, recipeAll)
	createRecipeNoItems := arbor.New("reject missing recipe items", testCreateRecipeNoItems, recipeAll)

	getRecipe := arbor.New("get recipe by id", testGetRecipe, createRecipeOk)
	getAllRecipes := arbor.New("get all recipes", testGetAllRecipe, createRecipeOk)

	disableRecipe := arbor.New("disable recipe", testDisableRecipe, getRecipe)

	placeOrderOK := arbor.New("post order", testCreateOrderOK, createRecipeOk, provision, createRecipeOk)
	placeOrderNoStock := arbor.New("rejects order with no stock", testCreateOrderWhenNotEnoughStock, createRecipeOk, provision, createRecipeOk)

	all, success := arbor.Run(createEmpty, getOne, getAll,
		noDuplicate, disable, getOneSP, getAllSP,
		getRecipe, createRecipeReject, createRecipeNoItems, getAllRecipes, disableRecipe,
		placeOrderOK, placeOrderNoStock,
	)

	t.Run("succeeds", func(t *testing.T) {
		assert.True(t, success)
	})

	t.Logf("\n%s\n", all)

	report := arbor.Marshal(all...)

	arbor.Upload(*arborURL, report)
}
