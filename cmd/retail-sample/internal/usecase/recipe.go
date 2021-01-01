package usecase

import (
	"context"
	"errors"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	domain "github.com/anatollupacescu/retail-sample/domain/retail/recipe"
)

func NewRecipe(ctx context.Context, book domain.Recipes, recipeDB recipeDB) Recipe {
	logger := log.Ctx(ctx).With().Str("layer", "usecase").Logger()

	return Recipe{
		ctx:      ctx,
		recipes:  book,
		recipeDB: recipeDB,
		logger:   &logger,
	}
}

type recipeDB interface {
	Add(domain.RecipeDTO) (domain.ID, error)
	Find(domain.Name) (*domain.RecipeDTO, error)
	Save(*domain.RecipeDTO) error
	Get(domain.ID) (domain.RecipeDTO, error)
}

type Recipe struct {
	logger   *zerolog.Logger
	recipes  domain.Recipes
	recipeDB recipeDB
	ctx      context.Context
}

type CreateRecipeDTO struct {
	Name        domain.Name
	Ingredients []domain.InventoryItem
}

func (o *Recipe) Create(dto CreateRecipeDTO) (recipe domain.RecipeDTO, err error) {
	id, err := o.recipes.Add(dto.Name, dto.Ingredients)

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

func (o *Recipe) Update(in UpdateStatusDTO) (domain.RecipeDTO, error) {
	recipeID := domain.ID(in.RecipeID)

	dto, err := o.recipeDB.Get(recipeID)

	if err != nil {
		return domain.RecipeDTO{}, err
	}

	recipe := domain.Recipe{
		ID:          dto.ID,
		Name:        dto.Name,
		Ingredients: dto.Ingredients,
		Enabled:     dto.Enabled,
		DB:          o.recipeDB,
	}

	switch in.Enabled {
	case true:
		err = recipe.Enable()
	default:
		err = recipe.Disable()
	}

	if err != nil {
		o.logger.Error().Err(err).Msg("call domain layer")
		return domain.RecipeDTO{}, err
	}

	dto.Enabled = recipe.Enabled

	return dto, nil
}
