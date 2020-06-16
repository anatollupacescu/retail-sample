package acceptance_test

import (
	"errors"
	"fmt"
	"math/rand"

	domain "github.com/anatollupacescu/retail-sample/internal/retail-domain/recipe"

	http "github.com/anatollupacescu/retail-sample/cmd/retail-sample-test"
	random "github.com/anatollupacescu/retail-sample/cmd/retail-sample-test"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/app/recipe"
)

func createRecipe(name string, items map[int]int) (domain.Recipe, error) {
	cl := http.Post("recipe")

	return recipe.Create(name, items, cl)
}

func testCreateRecipe() error {
	id := createRandomItem()

	qty := rand.Intn(100) + 1

	ingredients := map[int]int{
		id: qty,
	}

	recipeName := random.Name()

	r, err := createRecipe(recipeName, ingredients)

	if err != nil {
		return fmt.Errorf("could not create recipe: %v", err)
	}

	if r.Name != domain.Name(recipeName) {
		return fmt.Errorf("bad name")
	}

	return nil
}

func testCreateRecipeNoName() error {
	id := createRandomItem()

	qty := rand.Intn(100) + 1

	ingredients := map[int]int{
		id: qty,
	}

	_, err := createRecipe("", ingredients)

	if err == nil {
		return errors.New("expected error")
	}

	return nil
}

func testCreateRecipeNoItems() error {
	name := random.Name()

	_, err := createRecipe(name, nil)

	if err == nil {
		return errors.New("expected error")
	}

	return nil
}

func testGetRecipe() error {
	id := createRandomItem()

	qty := rand.Intn(100) + 1

	ingredients := map[int]int{
		id: qty,
	}

	recipeName := random.Name()

	r, err := createRecipe(recipeName, ingredients)

	if err != nil {
		return err
	}

	gcl := http.Get("recipe", int(r.ID))

	r, err = recipe.Get(gcl)

	if err != nil {
		return err
	}

	if r.ID == 0 {
		return errors.New("bad ID")
	}

	//TODO check ingredients

	return nil
}

func testGetAllRecipe() error {
	id := createRandomItem()

	ingredients := map[int]int{
		id: 4,
	}

	name := random.Name()

	_, _ = createRecipe(name, ingredients)

	gcl := http.List("recipe")

	all, err := recipe.GetAll(gcl)

	if err != nil {
		return err
	}

	if len(all) == 0 {
		return errors.New("no recipes")
	}

	return nil
}

func testDisableRecipe() error {
	id := createRandomItem()

	ingredients := map[int]int{
		id: 41,
	}

	name := random.Name()

	r, _ := createRecipe(name, ingredients)

	cl := http.Patch("recipe", int(r.ID))

	updated, err := recipe.Update(false, cl)

	if err != nil {
		return err
	}

	if updated.Enabled {
		return errors.New("expected to be disabled")
	}

	gcl := http.Get("recipe", int(r.ID))

	fetched, _ := recipe.Get(gcl)

	if fetched.ID != r.ID {
		return errors.New("not the same recipe")
	}

	if fetched.Enabled {
		return errors.New("expected to be disabled")
	}

	return nil
}
