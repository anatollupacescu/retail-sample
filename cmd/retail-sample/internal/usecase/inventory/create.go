package inventory

import (
	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
)

type CreateDTO struct {
	Name string
}

func (a *Inventory) Create(in CreateDTO) (item inventory.Item, err error) {
	a.logger.Info("create", "enter")

	var id int

	if id, err = a.inventory.Add(in.Name); err != nil {
		a.logger.Error("create", "call domain", err)
		return
	}

	if item, err = a.store.Get(id); err != nil {
		a.logger.Error("create", "retrieve new item", err)
		return
	}

	a.logger.Info("create", "success")

	return item, nil
}
