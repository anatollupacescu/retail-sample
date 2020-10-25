package stock

import (
	"errors"
	"strconv"

	"github.com/anatollupacescu/retail-sample/domain/retail/stock"
)

var ErrBadItemID = errors.New("could not parse ID")

func (o *Stock) Position(stockID string) (r stock.Position, err error) {
	o.logger.Info("quantity", "enter")

	var id int

	id, err = strconv.Atoi(stockID)

	if err != nil {
		o.logger.Error("quantity", "convert request ID", err)
		return stock.Position{}, ErrBadItemID
	}

	r, err = o.stock.Position(id)

	if err != nil {
		o.logger.Error("quantity", "call domain layer", err)
		return stock.Position{}, err
	}

	return r, nil
}
