package order

import "github.com/anatollupacescu/retail-sample/domain/retail/order"

func (o *UseCase) Create(recipeID, count int) (order.DTO, error) {
	id, err := o.orders.Create(recipeID, count)
	if err != nil {
		o.logger.Error().Err(err).Msg("call domain layer")
		return order.DTO{}, err
	}

	newOrder, err := o.orderDB.Get(id)
	if err != nil {
		o.logger.Error().Err(err).Msg("retrieve new order")
		return order.DTO{}, err
	}

	return newOrder, nil
}
