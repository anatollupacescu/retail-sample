package persistence

import (
	"context"
	"log"

	"github.com/anatollupacescu/retail-sample/internal/retail-domain/recipe"
)

type PgxRecipeStore struct {
	DB PgxDB
}

func (pr *PgxRecipeStore) Add(recipe.Recipe) (recipe.ID, error) {
	return recipe.ID(0), nil
}

func (pr *PgxRecipeStore) Get(id recipe.ID) recipe.Recipe {
	var name string
	err := pr.DB.QueryRow(context.Background(), "select name from recipe where id = $1", id).Scan(&name)

	if err != nil {
		log.Fatal(err)
	}

	return recipe.Recipe{}
}

func (pr *PgxRecipeStore) List() []recipe.Recipe {
	return nil
}
