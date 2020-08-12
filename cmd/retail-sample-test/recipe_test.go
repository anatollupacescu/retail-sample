package acceptance_test

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/anatollupacescu/arbortest/runner"

	domain "github.com/anatollupacescu/retail-sample/internal/retail-domain/recipe"

	http "github.com/anatollupacescu/retail-sample/cmd/retail-sample-test"
	random "github.com/anatollupacescu/retail-sample/cmd/retail-sample-test"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/app/recipe"
)

func createRecipe(name string, items map[int]int) (domain.Recipe, error) {
	cl := http.Post("recipe")

	return recipe.Create(name, items, cl)
}

// group:recipe after:inventory
func testCreateRecipe(t *runner.T) {
	createdName := random.Name()
	item, _ := createItem(createdName)

	qty := rand.Intn(100) + 1

	ingredients := map[int]int{
		item.ID: qty,
	}

	recipeName := random.Name()

	r, err := createRecipe(recipeName, ingredients)

	if err != nil {
		t.Error(fmt.Errorf("could not create recipe: %v", err))
		return
	}

	if r.Name != domain.Name(recipeName) {
		t.Error(fmt.Errorf("bad name"))
	}

	if len(r.Ingredients) != 1 {
		t.Error("expected one ingredient")
	}

	for _, i := range r.Ingredients {
		if i.ID != item.ID {
			t.Errorf("bad ingredient id, wanted %d got %d", item.ID, i.ID)
		}
		if i.Qty != qty {
			t.Errorf("bad qty, wanted %d got %d", qty, i.Qty)
		}
	}
}

// group:recipe
func testCreateRecipeNoName(t *runner.T) {
	createdName := random.Name()
	item, _ := createItem(createdName)

	id := item.ID

	qty := rand.Intn(100) + 1

	ingredients := map[int]int{
		id: qty,
	}

	_, err := createRecipe("", ingredients)

	if err == nil {
		t.Error(errors.New("expected error"))
	}
}

// group:recipe
func testCreateRecipeNoItems(t *runner.T) {
	name := random.Name()

	_, err := createRecipe(name, nil)

	if err == nil {
		t.Error(errors.New("expected error"))
	}
}

// group:recipe
func testGetRecipe(t *runner.T) {
	createdName := random.Name()
	item, _ := createItem(createdName)

	id := item.ID

	qty := rand.Intn(100) + 1

	ingredients := map[int]int{
		id: qty,
	}

	recipeName := random.Name()

	r, err := createRecipe(recipeName, ingredients)

	if err != nil {
		t.Error(err)
		return
	}

	gcl := http.Get("recipe", int(r.ID))

	r, err = recipe.Get(gcl)

	if err != nil {
		t.Error(err)
		return
	}

	if r.ID == 0 {
		t.Error(errors.New("bad ID"))
	}

	//TODO check ingredients
}

// group:recipe
func testGetAllRecipes(t *runner.T) {
	createdName := random.Name()
	item, _ := createItem(createdName)

	id := item.ID

	ingredients := map[int]int{
		id: 4,
	}

	name := random.Name()

	_, _ = createRecipe(name, ingredients)

	gcl := http.List("recipe")

	all, err := recipe.GetAll(gcl)

	if err != nil {
		t.Error(err)
		return
	}

	if len(all) == 0 {
		t.Error(errors.New("no recipes"))
	}
}

// group:recipe
func testDisableRecipe(t *runner.T) {
	createdName := random.Name()
	item, _ := createItem(createdName)

	id := item.ID

	ingredients := map[int]int{
		id: 41,
	}

	name := random.Name()

	r, _ := createRecipe(name, ingredients)

	cl := http.Patch("recipe", int(r.ID))

	updated, err := recipe.Update(false, cl)

	if err != nil {
		t.Error(err)
		return
	}

	if updated.Enabled {
		t.Error(errors.New("expected to be disabled"))
		return
	}

	gcl := http.Get("recipe", int(r.ID))

	fetched, _ := recipe.Get(gcl)

	if fetched.ID != r.ID {
		t.Error(errors.New("not the same recipe"))
		return
	}

	if fetched.Enabled {
		t.Error(errors.New("expected to be disabled"))
	}
}
