package stock

import (
	"context"

	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
	"github.com/anatollupacescu/retail-sample/domain/retail/stock"
)

type logger interface {
	Error(string, string, error)
	Info(string, string)
}

func New(ctx context.Context, stock stock.Stock, provisionLog stock.ProvisionLog,
	stockDB stock.Store,
	inventoryDB inventory.Store,
	log logger) Stock {
	return Stock{
		ctx:          ctx,
		stock:        stock,
		provisionLog: provisionLog,
		stockDB:      stockDB,
		inventoryDB:  inventoryDB,
		logger:       log,
	}
}

type Stock struct {
	logger       logger
	stock        stock.Stock
	provisionLog stock.ProvisionLog
	stockDB      stock.Store
	inventoryDB  inventory.Store
	ctx          context.Context
}
