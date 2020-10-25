package inventory

import (
	"strconv"

	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
)

func (a *Inventory) UpdateStatus(rid string, status bool) (inventory.Item, error) {
	a.logger.Info("update status", "begin")

	id, err := strconv.Atoi(rid)

	if err != nil {
		a.logger.Error("update status", "parse item id", err)
		return inventory.Item{}, err
	}

	var item inventory.Item

	if item, err = a.inventory.UpdateStatus(id, status); err != nil {
		a.logger.Error("update status", "call domain", err)
		return inventory.Item{}, err
	}

	a.logger.Info("update status", "success")

	return item, nil
}
