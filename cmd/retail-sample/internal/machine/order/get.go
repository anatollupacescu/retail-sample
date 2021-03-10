package order

import (
	"strconv"

	"github.com/pkg/errors"

	usecase "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/machine"
	domain "github.com/anatollupacescu/retail-sample/domain/retail/order"
)

func (o *UseCase) GetByID(orderID string) (domain.DTO, error) {
	var err error

	defer func() {
		if err != nil {
			o.logger.Error().Str("action", "get by id").Err(err).Send()
		}
	}()

	id, err := strconv.Atoi(orderID)

	if err != nil {
		return domain.DTO{}, errors.Wrapf(usecase.ErrBadRequest, "parse order ID: %v", orderID)
	}

	order, err := o.orderDB.Get(id)
	switch err {
	case nil:
	case domain.ErrOrderNotFound:
		return domain.DTO{}, errors.Wrapf(usecase.ErrNotFound, "get order with id: %v", id)
	default:
		return domain.DTO{}, err
	}

	return order, nil
}

func (o *UseCase) GetAll() ([]domain.DTO, error) {
	orders, err := o.orderDB.List()
	if err != nil {
		o.logger.Error().Str("action", "get all").Err(err).Send()
		return nil, err
	}

	return orders, nil
}
