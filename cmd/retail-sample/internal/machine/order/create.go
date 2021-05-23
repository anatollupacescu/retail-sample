package order

import (
	"github.com/pkg/errors"

	usecase "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/machine"
	"github.com/anatollupacescu/retail-sample/domain/retail/order"
	"github.com/anatollupacescu/retail-sample/domain/retail/recipe"
	"github.com/anatollupacescu/retail-sample/domain/retail/stock"
)

func (o *UseCase) Create(recipeID, count int) (order.DTO, error) {
	var err error

	defer func() {
		if err != nil {
			o.logger.Error().Str("action", "create order").Err(err).Send()
		}
	}()

	id, err := o.orders.Create(recipeID, count)

	switch err {
	case nil:
	case recipe.ErrNotFound:
		return order.DTO{}, errors.Wrapf(usecase.ErrNotFound, "create order for recipe %d: %v", recipeID, err)
	case
		order.ErrInvalidRecipe,
		order.ErrInvalidQuantity,
		stock.ErrInvalidExtractQuantity,
		stock.ErrNotEnoughStock:
		return order.DTO{}, errors.Wrapf(usecase.ErrBadRequest, "create order for recipe %d: %v", recipeID, err)
	default:
		return order.DTO{}, err
	}

	o.logger.Info().Int("id", id).Msg("successfully created order")

	newOrder, err := o.orderDB.Get(id)

	if err != nil {
		return order.DTO{}, err
	}

	return newOrder, nil
}
