package inventory

import (
	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
)

type UpdateStatusDTO struct {
	ID      int
	Enabled bool
}

func (a *Inventory) UpdateStatus(in UpdateStatusDTO) (item inventory.Item, err error) {
	a.logger.Info("update status", "begin")

	if item, err = a.inventory.UpdateStatus(in.ID, in.Enabled); err != nil {
		a.logger.Error("update status", "call domain", err)
		return inventory.Item{}, err
	}

	a.logger.Info("update status", "success")

	return item, nil
}
