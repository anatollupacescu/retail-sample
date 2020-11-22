package usecase

import (
	"context"
	"errors"

	"github.com/anatollupacescu/retail-sample/domain/retail/recipe"
)

func NewRecipe(ctx context.Context, book recipe.Book, recipeDB recipeDB, log logger) Recipe {
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

type CreateRecipeDTO struct {
	Name        recipe.Name
	Ingredients []recipe.Ingredient
}

func (o *Recipe) Create(dto CreateRecipeDTO) (recipe recipe.Recipe, err error) {
	o.logger.Info("get all", "enter")

	id, err := o.book.Add(dto.Name, dto.Ingredients)

	if err != nil {
		o.logger.Error("get all", "call domain layer", err)
		return
	}

	recipe, err = o.book.DB.Get(id)

	if err != nil {
		o.logger.Error("get all", "call domain layer to retrieve the newly created recipe", err)
		return
	}

	o.logger.Info("get all", "success")

	return
}

var ErrBadItemID = errors.New("could not parse ID")

type UpdateStatusDTO struct {
	RecipeID int
	Enabled  bool
}

func (o *Recipe) Update(in UpdateStatusDTO) (recipe.Recipe, error) {
	o.logger.Info("update status", "begin")

	recipeID := recipe.ID(in.RecipeID)

	err := o.book.SetStatus(recipeID, in.Enabled)

	if err != nil {
		o.logger.Error("update status", "call domain", err)
		return recipe.Recipe{}, err
	}

	rec, err := o.recipeDB.Get(recipeID)

	if err != nil {
		o.logger.Error("update status", "call domain to retrieve updated entity", err)
		return recipe.Recipe{}, err
	}

	o.logger.Info("update status", "success")

	return rec, nil
}
