package acceptance_test

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/anatollupacescu/retail-sample/internal/arbor"
)

var arborURL = flag.String("arborURL", "", "graph server URL")

func TestAcceptance(t *testing.T) {
	createOk := arbor.New("can create", testCreate)
	createEmpty := arbor.New("rejects empty name", testCreateWithEmptyName, createOk)

	getOne := arbor.New("get one", testGetOne, createOk)
	getAll := arbor.New("get all", testGetAll, createOk)
	noDuplicate := arbor.New("no duplicate", testDuplicate, createOk)
	disable := arbor.New("disable", testDisable, createOk)

	//stock
	provision := arbor.New("provision stock", testProvision, createOk)
	getOneSP := arbor.New("get single stock position", testGetStockPos, provision)
	getAllSP := arbor.New("get all stock positions", testGetAllStockPos, provision)

	//recipe
	recipeAll := arbor.Alias("create recipe ingredient", createOk)

	createRecipeOk := arbor.New("can create recipe", testCreateRecipe, recipeAll)
	createRecipeReject := arbor.New("reject empty name", testCreateRecipeNoName, recipeAll)
	createRecipeNoItems := arbor.New("reject missing items", testCreateRecipeNoItems, recipeAll)

	getRecipe := arbor.New("get recipe by id", testGetRecipe, createRecipeOk)
	getAllRecipes := arbor.New("get all recipes", testGetAllRecipe, createRecipeOk)

	disableRecipe := arbor.New("disable recipe", testDisableRecipe, getRecipe)

	all, success := arbor.Run(createEmpty, getOne, getAll,
		noDuplicate, disable, getOneSP, getAllSP,
		getRecipe, createRecipeReject, createRecipeNoItems, getAllRecipes, disableRecipe)

	t.Run("succeeds", func(t *testing.T) {
		assert.Equal(t, true, success)
	})

	t.Logf("\n%s\n", all)

	report := arbor.Marshal(all...)

	arbor.Upload(*arborURL, report)
}
