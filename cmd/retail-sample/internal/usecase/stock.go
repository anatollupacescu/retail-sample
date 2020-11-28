package usecase

import (
	"context"

	"github.com/anatollupacescu/retail-sample/domain/retail/stock"
)

func NewStock(ctx context.Context, stock stock.Stock,
	stockDB stockDB, logDB logDB, inventoryDB inventoryDB,
	log logger) Stock {
	return Stock{
		ctx:         ctx,
		stock:       stock,
		stockDB:     stockDB,
		logDB:       logDB,
		inventoryDB: inventoryDB,
		logger:      log,
	}
}

type stockDB interface {
	Quantity(id int) (int, error)
}

type logDB interface {
	Add(id, qty int) (int, error)
}

type Stock struct {
	logger      logger
	stock       stock.Stock
	stockDB     stockDB
	inventoryDB inventoryDB
	logDB       logDB
	ctx         context.Context
}

type ProvisionDTO struct {
	InventoryItemID int
	Qty             int
}

type Position struct {
	ID   int
	Name string
	Qty  int
}

func (o *Stock) Provision(dto ProvisionDTO) (Position, error) {
	o.logger.Info("provision", "enter")

	err := o.stock.Provision(dto.InventoryItemID, dto.Qty)
	if err != nil {
		o.logger.Error("provision", "call domain layer", err)
		return Position{}, err
	}

	_, err = o.logDB.Add(dto.InventoryItemID, dto.Qty)
	if err != nil {
		o.logger.Error("provision", "call domain layer", err)
		return Position{}, err
	}

	qty, err := o.stockDB.Quantity(dto.InventoryItemID)
	if err != nil {
		o.logger.Error("provision", "call domain layer to retrieve quantity", err)
		return Position{}, err
	}

	item, err := o.inventoryDB.Get(dto.InventoryItemID)
	if err != nil {
		o.logger.Error("provision", "call domain layer to retrieve stock position", err)
		return Position{}, err
	}

	o.logger.Error("provision", "success", err)

	pos := Position{
		ID:   item.ID,
		Name: item.Name,
		Qty:  qty,
	}

	return pos, nil
}
