package usecase

import (
	"context"
	"errors"

	"github.com/anatollupacescu/retail-sample/domain/retail/recipe"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func NewRecipe(ctx context.Context, book recipe.Collection, recipeDB recipeDB) Recipe {
	logger := log.Ctx(ctx).With().Str("layer", "usecase").Logger()

	return Recipe{
		ctx:      ctx,
		book:     book,
		recipeDB: recipeDB,
		logger:   &logger,
	}
}

type recipeDB interface {
	Get(recipe.ID) (recipe.Recipe, error)
}

type Recipe struct {
	logger   *zerolog.Logger
	book     recipe.Collection
	recipeDB recipeDB
	ctx      context.Context
}

type CreateRecipeDTO struct {
	Name        recipe.Name
	Ingredients []recipe.Ingredient
}

func (o *Recipe) Create(dto CreateRecipeDTO) (recipe recipe.Recipe, err error) {
	id, err := o.book.Add(dto.Name, dto.Ingredients)

	if err != nil {
		o.logger.Error().Err(err).Msg("call domain layer")
		return
	}

	recipe, err = o.recipeDB.Get(id)

	if err != nil {
		o.logger.Error().Err(err).Msg("retrieve the newly created recipe")
		return
	}

	return
}

var ErrBadItemID = errors.New("could not parse ID")

type UpdateStatusDTO struct {
	RecipeID int
	Enabled  bool
}

func (o *Recipe) Update(in UpdateStatusDTO) (recipe.Recipe, error) {
	recipeID := recipe.ID(in.RecipeID)

	found, err := o.recipeDB.Get(recipeID)

	if err != nil {
		return recipe.Recipe{}, err
	}

	switch in.Enabled {
	case true:
		err = found.Enable()
	default:
		err = found.Disable()
	}

	if err != nil {
		o.logger.Error().Err(err).Msg("call domain layer")
		return recipe.Recipe{}, err
	}

	return found, nil
}
