package stock

import "github.com/anatollupacescu/retail-sample/domain/retail/stock"

func (o *Stock) ProvisionLog() ([]stock.ProvisionEntry, error) {
	o.logger.Info("provision log", "enter")

	all, err := o.stock.GetAllProvisions()
	if err != nil {
		o.logger.Error("provision log", "call domain layer", err)
		return nil, err
	}

	o.logger.Error("provision log", "success", err)

	return all, nil
}
