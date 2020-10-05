package recipe

import (
	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/middleware"
	"github.com/anatollupacescu/retail-sample/domain/retail-sample/recipe"
)

type wrapper struct {
	middleware.Middleware
}

func (w wrapper) setStatus(id int, enabled bool) (re recipe.Recipe, err error) {
	return re, w.Exec("disable recipe", func(provider middleware.PersistenceProvider) error {
		r := provider.RecipeBook()

		re, err = r.SetStatus(id, enabled)

		return err
	})
}

func (w wrapper) create(recipeName recipe.Name, recipeIngredients []recipe.Ingredient) (re recipe.Recipe, err error) {
	return re, w.Exec("add recipe", func(provider middleware.PersistenceProvider) error {
		var recipeID recipe.ID

		r := provider.RecipeBook()

		recipeID, err = r.Add(recipeName, recipeIngredients)

		if err != nil {
			return err
		}

		re, err = r.Get(recipeID)

		return err
	})
}

func (w wrapper) get(recipeID recipe.ID) (out recipe.Recipe, err error) {
	return out, w.Exec("get recipe", func(provider middleware.PersistenceProvider) error {
		r := provider.RecipeBook()

		out, err = r.Get(recipeID)

		return err
	})
}

func (w wrapper) getAll() (recipes []recipe.Recipe, err error) {
	return recipes, w.Exec("get recipe", func(provider middleware.PersistenceProvider) error {
		r := provider.RecipeBook()

		recipes, err = r.List()

		return err
	})
}
