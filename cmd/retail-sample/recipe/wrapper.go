package recipe

import (
	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/types"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/recipe"
)

type wrapper struct {
	loggerFactory              types.LoggerFactory
	persistenceProviderFactory types.PersistenceProviderFactory
}

func (w *wrapper) exec(methodName string, f func(recipe.Book) error) {
	logger := w.loggerFactory()

	logger.Log("msg", "enter", "method", methodName)
	defer logger.Log("msg", "exit", "method", methodName)

	provider := w.persistenceProviderFactory.New()

	recipes := provider.RecipeBook()

	err := f(recipes)

	if err != nil {
		logger.Log("msg", "rollback")
		w.persistenceProviderFactory.Rollback(provider)
		return
	}

	logger.Log("msg", "commit")
	w.persistenceProviderFactory.Commit(provider)
}

func (w wrapper) setStatus(id int, enabled bool) (re recipe.Recipe, err error) {
	w.exec("disable recipe", func(r recipe.Book) error {
		re, err = r.SetStatus(id, enabled)

		return err
	})

	return
}

func (w wrapper) create(recipeName recipe.Name, recipeIngredients []recipe.Ingredient) (re recipe.Recipe, err error) {
	w.exec("add recipe", func(r recipe.Book) error {
		var recipeID recipe.ID

		recipeID, err = r.Add(recipeName, recipeIngredients)

		if err != nil {
			return err
		}

		re, err = r.Get(recipeID)

		return err
	})

	return
}

func (w wrapper) get(recipeID recipe.ID) (out recipe.Recipe, err error) {
	w.exec("get recipe", func(r recipe.Book) error {
		out, err = r.Get(recipeID)

		return err
	})

	return
}

func (w wrapper) getAll() (recipes []recipe.Recipe, err error) {
	w.exec("get recipe", func(r recipe.Book) error {
		recipes, err = r.List()

		return err
	})

	return
}
