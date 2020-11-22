package usecase

import (
	"context"

	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
)

type inventoryDB interface {
	Get(int) (inventory.Item, error)
}

func NewInventory(ctx context.Context, inventory inventory.Inventory, db inventoryDB, log logger) Inventory {
	return Inventory{
		ctx:         ctx,
		inventory:   inventory,
		inventoryDB: db,
		logger:      log,
	}
}

type Inventory struct {
	logger      logger
	inventory   inventory.Inventory
	inventoryDB inventoryDB
	ctx         context.Context
}

type CreateInventoryItemDTO struct {
	Name string
}

func (a *Inventory) Create(in CreateInventoryItemDTO) (item inventory.Item, err error) {
	a.logger.Info("create", "enter")

	var id int

	if id, err = a.inventory.Add(in.Name); err != nil {
		a.logger.Error("create", "call domain", err)
		return
	}

	if item, err = a.inventoryDB.Get(id); err != nil {
		a.logger.Error("create", "retrieve new item", err)
		return
	}

	a.logger.Info("create", "success")

	return item, nil
}

type UpdateInventoryItemStatusDTO struct {
	ID      int
	Enabled bool
}

func (a *Inventory) UpdateStatus(in UpdateInventoryItemStatusDTO) (item inventory.Item, err error) {
	a.logger.Info("update status", "begin")

	if err = a.inventory.UpdateStatus(in.ID, in.Enabled); err != nil {
		a.logger.Error("update status", "call domain", err)
		return inventory.Item{}, err
	}

	item, err = a.inventory.DB.Get(in.ID)

	if err != nil {
		a.logger.Error("update status", "fetch updated item", err)
		return inventory.Item{}, err
	}

	a.logger.Info("update status", "success")

	return item, nil
}
