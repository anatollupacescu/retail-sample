package stock

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	pg "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/persistence/postgres"
	"github.com/anatollupacescu/retail-sample/domain/retail/stock"
)

func New(ctx context.Context, t pg.TX) UseCase {
	logger := log.Ctx(ctx).With().Str("layer", "usecase").Logger()

	stockDB := &pg.StockPgxStore{DB: t.Tx}
	logDB := &pg.PgxProvisionLog{DB: t.Tx}
	inventoryDB := &pg.InventoryPgxStore{DB: t.Tx}

	return UseCase{
		ctx:         ctx,
		stockDB:     stockDB,
		logDB:       logDB,
		inventoryDB: inventoryDB,
		logger:      &logger,
	}
}

type UseCase struct {
	logger      *zerolog.Logger
	stockDB     *pg.StockPgxStore
	inventoryDB *pg.InventoryPgxStore
	logDB       *pg.PgxProvisionLog
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

func (o *UseCase) Provision(dto ProvisionDTO) (Position, error) {
	stockPos, err := o.stockDB.Get(dto.InventoryItemID)

	switch err {
	case nil, pg.ErrStockItemNotFound: //continue
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
