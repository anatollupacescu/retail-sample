package stock

import (
	"github.com/pkg/errors"
)

type (
	Position struct {
		Validator   inventoryValidator
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

	inventoryValidator interface {
		Valid(id int) (bool, error)
	}
)

var (
	ErrInvalidProvisionQuantity = errors.New("invalid provision quantity")
	ErrPositionNotFound         = errors.New("stock position not found")
	ErrItemInvalid              = errors.New("item invalid")
)

func (p *Position) Provision(qty int) error {
	if qty <= 0 {
		return ErrInvalidProvisionQuantity
	}

	valid, err := p.Validator.Valid(p.InventoryID)

	if err != nil {
		return err
	}

	if !valid {
		return ErrItemInvalid
	}

	current := p.Qty + qty

	dto := PositionDTO{
		InventoryID: p.InventoryID,
		Qty:         current,
	}

	err = p.DB.Save(dto)

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

func (p *Position) extract(qty int) error {
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
