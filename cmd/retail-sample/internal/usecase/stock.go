package usecase

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/anatollupacescu/retail-sample/domain/retail/stock"
)

func NewStock(ctx context.Context, stock stock.Stock, stockDB stockDB, logDB logDB, inventoryDB inventoryDB) Stock {
	logger := log.Ctx(ctx).With().Str("layer", "usecase").Logger()

	return Stock{
		ctx:         ctx,
		stock:       stock,
		stockDB:     stockDB,
		logDB:       logDB,
		inventoryDB: inventoryDB,
		logger:      &logger,
	}
}

type stockDB interface {
	Quantity(id int) (int, error)
}

type logDB interface {
	Add(id, qty int) (int, error)
}

type Stock struct {
	logger      *zerolog.Logger
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
	err := o.stock.Provision(dto.InventoryItemID, dto.Qty)
	if err != nil {
		o.logger.Error().Err(err).Msg("call domain layer")
		return Position{}, err
	}

	_, err = o.logDB.Add(dto.InventoryItemID, dto.Qty)
	if err != nil {
		o.logger.Error().Err(err).Msg("call domain layer")
		return Position{}, err
	}

	qty, err := o.stockDB.Quantity(dto.InventoryItemID)
	if err != nil {
		o.logger.Error().Err(err).Msg("call domain layer to retrieve quantity")
		return Position{}, err
	}

	item, err := o.inventoryDB.Get(dto.InventoryItemID)
	if err != nil {
		o.logger.Error().Err(err).Msg("call domain layer to retrieve stock position")
		return Position{}, err
	}

	pos := Position{
		ID:   item.ID,
		Name: item.Name,
		Qty:  qty,
	}

	return pos, nil
}
