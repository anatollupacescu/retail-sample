package inventory

import (
	"context"

	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
)

type logger interface {
	Error(string, string, error)
	Info(string, string)
}

type inventoryDB interface {
	Get(int) (inventory.Item, error)
}

func New(ctx context.Context, inventory inventory.Inventory, db inventoryDB, log logger) Inventory {
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
