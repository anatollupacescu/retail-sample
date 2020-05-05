package recipe

import (
	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/types"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/recipe"
)

type wrapper struct {
	loggerFactory              types.LoggerFactory
	persistenceProviderFactory types.PersistenceProviderFactory

	provider types.PersistenceProvider
	logger   types.Logger
}

func (w *wrapper) exec(methodName string, f func() error) {
	logger := w.loggerFactory()

	logger.Log("msg", "enter", "method", methodName)
	defer logger.Log("msg", "exit", "method", methodName)

	w.provider = w.persistenceProviderFactory.New()

	err := f()

	if err != nil {
		logger.Log("msg", "rollback")
		w.persistenceProviderFactory.Rollback(w.provider)
		return
	}

	logger.Log("msg", "commit")
	w.persistenceProviderFactory.Commit(w.provider)
}

func (w wrapper) Add(recipeName recipe.Name, recipeIngredients []recipe.Ingredient) (recipeID recipe.ID, err error) {
	w.exec("add recipe", func() error {
		recipes := w.provider.RecipeBook()
		recipeID, err = recipes.Add(recipeName, recipeIngredients)

		return err
	})

	return
}

func (w wrapper) Get(recipeID recipe.ID) (recipe recipe.Recipe, err error) {
	w.exec("get recipe", func() error {
		recipes := w.provider.RecipeBook()
		recipe, err = recipes.Get(recipeID)

		return err
	})

	return
}

func (w wrapper) List() (rcps []recipe.Recipe, err error) {
	w.exec("get recipe", func() error {
		recipes := w.provider.RecipeBook()
		rcps, err = recipes.List()

		return err
	})

	return
}
