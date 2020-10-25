package stock

import (
	"strconv"

	"github.com/anatollupacescu/retail-sample/domain/retail/stock"
)

func (o *Stock) Provision(reqID string, qty int) (stock.Position, error) {
	o.logger.Info("provision", "enter")

	id, err := strconv.Atoi(reqID)
	if err != nil {
		o.logger.Error("provision", "convert request ID", err)
		return stock.Position{}, ErrBadItemID
	}

	provisionID, err := o.stock.Provision(id, qty)
	if err != nil {
		o.logger.Error("provision", "call domain layer", err)
		return stock.Position{}, err
	}

	logEntry, err := o.stock.GetProvision(provisionID)
	if err != nil {
		o.logger.Error("provision", "call domain layer to retrieve provision record", err)
		return stock.Position{}, err
	}

	pos, err := o.stock.Position(logEntry.ID)
	if err != nil {
		o.logger.Error("provision", "call domain layer to retrieve stock position", err)
		return stock.Position{}, err
	}

	o.logger.Error("provision", "success", err)

	return pos, nil
}
