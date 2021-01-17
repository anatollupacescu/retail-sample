package stock

import (
	"github.com/pkg/errors"
)

type (
	Position struct {
		InventoryID int
		Qty         int

		DB db
	}

	PositionDTO struct {
		InventoryID int
		Qty         int
	}

	db interface {
		Save(PositionDTO) error
		Get(int) (PositionDTO, error)
	}
)

var (
	ErrInvalidProvisionQuantity = errors.New("invalid provision quantity")
	ErrPositionNotFound         = errors.New("stock position not found")
)

func (p *Position) Provision(qty int) error {
	if qty <= 0 {
		return ErrInvalidProvisionQuantity
	}

	current := p.Qty + qty

	dto := PositionDTO{
		InventoryID: p.InventoryID,
		Qty:         current,
	}

	err := p.DB.Save(dto)

	if err != nil {
		return err
	}

	p.Qty = dto.Qty

	return nil
}

var (
	ErrNotEnoughStock         = errors.New("not enough stock")
	ErrInvalidExtractQuantity = errors.New("invalid extract quantity")
)

func (p *Position) Extract(qty int) error {
	if qty <= 0 {
		return ErrInvalidExtractQuantity
	}

	if qty > p.Qty {
		return ErrNotEnoughStock
	}

	current := p.Qty - qty

	dto := PositionDTO{
		InventoryID: p.InventoryID,
		Qty:         current,
	}

	err := p.DB.Save(dto)

	if err != nil {
		return err
	}

	p.Qty = dto.Qty

	return nil
}
