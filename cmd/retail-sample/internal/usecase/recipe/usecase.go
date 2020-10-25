package recipe

import (
	"context"

	"github.com/anatollupacescu/retail-sample/domain/retail/recipe"
)

type logger interface {
	Error(string, string, error)
	Info(string, string)
}

func New(ctx context.Context, book recipe.Book, log logger) Recipe {
	return Recipe{
		ctx:    ctx,
		book:   book,
		logger: log,
	}
}

type Recipe struct {
	logger logger
	book   recipe.Book
	ctx    context.Context
}
