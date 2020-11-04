package inventory

import (
	"context"

	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
)

type logger interface {
	Error(string, string, error)
	Info(string, string)
}

func New(ctx context.Context, inventory inventory.Inventory, store store, log logger) Inventory {
	return Inventory{
		ctx:       ctx,
		inventory: inventory,
		logger:    log,
		store:     store,
	}
}

type store interface {
	Get(int) (inventory.Item, error)
	List() ([]inventory.Item, error)
}

type Inventory struct {
	logger    logger
	inventory inventory.Inventory
	store     store
	ctx       context.Context
}
