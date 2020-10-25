package inventory

import (
	"errors"
	"strconv"

	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
)

var ErrBadItemID = errors.New("could not parse ID")

func (a *Inventory) GetByID(itemID string) (item inventory.Item, err error) {
	a.logger.Info("get by id", "enter")

	var id int

	id, err = strconv.Atoi(itemID)

	if err != nil {
		a.logger.Error("get by id", "convert request ID", err)
		return item, ErrBadItemID
	}

	item, err = a.inventory.Get(id)

	if err != nil {
		a.logger.Error("get by id", "call domain layer", err)
		return
	}

	a.logger.Info("get by id", "success")

	return item, nil
}

func (a *Inventory) GetAll() (items []inventory.Item, err error) {
	a.logger.Info("get all", "enter")

	items, err = a.inventory.List()

	if err != nil {
		a.logger.Error("get all", "call domain layer", err)
		return
	}

	a.logger.Info("get all", "success")

	return items, nil
}
