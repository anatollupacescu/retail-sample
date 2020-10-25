package inventory

import (
	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
)

func (a *Inventory) Create(name string) (item inventory.Item, err error) {
	a.logger.Info("create", "enter")

	var id int

	if id, err = a.inventory.Add(name); err != nil {
		a.logger.Error("create", "call domain", err)
		return
	}

	if item, err = a.inventory.Get(id); err != nil {
		a.logger.Error("create", "retrieve new item", err)
		return
	}

	a.logger.Info("create", "success")

	return item, nil
}
