package stock

import "github.com/anatollupacescu/retail-sample/domain/retail/stock"

func (o *Stock) CurrentStock() (stocks []stock.Position, err error) {
	o.logger.Info("get all", "enter")

	stocks, err = o.stock.CurrentStock()

	if err != nil {
		o.logger.Error("get all", "call domain layer", err)

		return
	}

	o.logger.Info("get all", "success")

	return
}
