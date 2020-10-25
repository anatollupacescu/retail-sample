package order

import (
	"errors"
	"strconv"

	"github.com/anatollupacescu/retail-sample/domain/retail/order"
)

var ErrBadItemID = errors.New("could not parse ID")

func (o *Order) GetByID(itemID string) (ordr order.Order, err error) {
	o.logger.Info("get by id", "enter")

	var id int

	id, err = strconv.Atoi(itemID)

	if err != nil {
		o.logger.Error("get by id", "convert request ID", err)

		return ordr, ErrBadItemID
	}

	ordr, err = o.orders.Get(order.ID(id))

	if err != nil {
		o.logger.Error("get by id", "call domain layer", err)

		return
	}

	o.logger.Info("get by id", "success")

	return ordr, nil
}
