package stock

import (
	"github.com/anatollupacescu/retail-sample/domain/retail/stock"
)

type UpdateDTO struct {
	ReqID int
	Qty   int
}

func (o *Stock) Provision(dto UpdateDTO) (stock.Position, error) {
	o.logger.Info("provision", "enter")

	provisionID, err := o.stock.Provision(dto.ReqID, dto.Qty)
	if err != nil {
		o.logger.Error("provision", "call domain layer", err)
		return stock.Position{}, err
	}

	logEntry, err := o.provisionLog.Get(provisionID)
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
