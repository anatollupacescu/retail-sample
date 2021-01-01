package usecase

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	persistence "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/persistence/postgres"
	"github.com/anatollupacescu/retail-sample/domain/retail/stock"
)

func NewStock(ctx context.Context, stockDB stockDB, logDB logDB, inventoryDB inventoryDB) Stock {
	logger := log.Ctx(ctx).With().Str("layer", "usecase").Logger()

	return Stock{
		ctx:         ctx,
		stockDB:     stockDB,
		logDB:       logDB,
		inventoryDB: inventoryDB,
		logger:      &logger,
	}
}

type stockDB interface {
	Get(int) (stock.PositionDTO, error)
	Save(stock.PositionDTO) error
}

type logDB interface {
	Add(id, qty int) (int, error)
}

type Stock struct {
	logger      *zerolog.Logger
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
	stockPos, err := o.stockDB.Get(dto.InventoryItemID)

	switch err {
	case nil, persistence.ErrStockItemNotFound: //continue
	default:
		return Position{}, err
	}

	pos := stock.Position{
		Qty:         stockPos.Qty,
		InventoryID: dto.InventoryItemID,
		DB:          o.stockDB,
	}

	err = pos.Provision(dto.Qty)
	if err != nil {
		o.logger.Error().Err(err).Msg("call domain layer")
		return Position{}, err
	}

	_, err = o.logDB.Add(dto.InventoryItemID, dto.Qty)
	if err != nil {
		o.logger.Error().Err(err).Msg("add log entry")
		return Position{}, err
	}

	item, err := o.inventoryDB.Get(dto.InventoryItemID)
	if err != nil {
		o.logger.Error().Err(err).Msg("get inventory item name")
		return Position{}, err
	}

	result := Position{
		ID:   item.ID,
		Name: item.Name,
		Qty:  pos.Qty,
	}

	return result, nil
}
