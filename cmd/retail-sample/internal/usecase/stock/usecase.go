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

func New(ctx context.Context, stock stock.Stock,
	stockDB stockDB, inventoryDB inventoryDB,
	log logger) Stock {
	return Stock{
		ctx:         ctx,
		stock:       stock,
		stockDB:     stockDB,
		inventoryDB: inventoryDB,
		logger:      log,
	}
}

type (
	stockDB interface {
		Quantity(id int) (int, error)
	}
	inventoryDB interface {
		Get(id int) (inventory.Item, error)
	}
)

type Stock struct {
	logger      logger
	stock       stock.Stock
	stockDB     stockDB
	inventoryDB inventoryDB
	ctx         context.Context
}
