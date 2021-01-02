package recipe

import (
	"context"
	"errors"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	pg "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/persistence/postgres"
	domain "github.com/anatollupacescu/retail-sample/domain/retail/recipe"
)

func New(ctx context.Context, t pg.TX) UseCase {
	logger := log.Ctx(ctx).With().Str("layer", "usecase").Logger()

	recipeDB := &pg.RecipePgxStore{DB: t.Tx}
	inventoryDB := &pg.InventoryPgxStore{DB: t.Tx}

	recipes := domain.Recipes{
		DB:        recipeDB,
		Inventory: &validator{Inventory: inventoryDB},
	}

	return UseCase{
		ctx:      ctx,
		recipes:  recipes,
		recipeDB: recipeDB,
		logger:   &logger,
	}
}

type UseCase struct {
	ctx      context.Context
	logger   *zerolog.Logger
	recipes  domain.Recipes
	recipeDB *pg.RecipePgxStore
}

type CreateRecipeDTO struct {
	Name        domain.Name
	Ingredients []domain.InventoryItem
}

func (o *UseCase) Create(dto CreateRecipeDTO) (recipe domain.RecipeDTO, err error) {
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

func (o *UseCase) UpdateStatus(in UpdateStatusDTO) (domain.RecipeDTO, error) {
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
