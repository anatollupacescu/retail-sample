package inventory

import (
	"context"

	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
)

type logger interface {
	Error(string, string, error)
	Info(string, string)
}

func New(ctx context.Context, inventory inventory.Inventory, log logger) Inventory {
	return Inventory{
		ctx:       ctx,
		inventory: inventory,
		logger:    log,
	}
}

type Inventory struct {
	logger    logger
	inventory inventory.Inventory
	ctx       context.Context
}
