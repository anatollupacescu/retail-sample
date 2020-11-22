package recipe

import (
	"context"

	"github.com/anatollupacescu/retail-sample/domain/retail/recipe"
)

type logger interface {
	Error(string, string, error)
	Info(string, string)
}

func New(ctx context.Context, book recipe.Book, recipeDB recipeDB, log logger) Recipe {
	return Recipe{
		ctx:      ctx,
		book:     book,
		recipeDB: recipeDB,
		logger:   log,
	}
}

type recipeDB interface {
	Get(recipe.ID) (recipe.Recipe, error)
}

type Recipe struct {
	logger   logger
	book     recipe.Book
	recipeDB recipeDB
	ctx      context.Context
}
