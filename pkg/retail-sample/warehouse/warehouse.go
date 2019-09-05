//go:generate mockgen -source=warehouse.go -package mocks -destination mocks/warehouse.go
package warehouse

import "github.com/anatollupacescu/retail-sample/pkg/retail-sample/itemtype"

type (
	Store interface {
		Add(uint64, int)
		Update(uint64, int) error
		Get(uint64) (int, error)
	}

	Repository struct {
		ItemStore          Store
		ItemTypeRepository itemtype.Repository
	}
)

const (
	ErrItemTypeNotFound = warehouseError("Item type with given id was not found")
	ErrItemNotFound     = warehouseError("No such item stored in the warehouse")
	ErrUpdate           = warehouseError("Could not update quantity")
)

var (
	zeroItemTypeValue = ""
)

func (r *Repository) Add(id uint64, qty int) error {
	itemType := r.ItemTypeRepository.Get(id)

	if itemType == zeroItemTypeValue {
		return ErrItemTypeNotFound
	}

	if _, err := r.ItemStore.Get(id); err != nil {
		r.ItemStore.Add(id, qty)
		return nil
	}

	if err := r.ItemStore.Update(id, qty); err != nil {
		return ErrUpdate
	}

	return nil
}

type warehouseError string

func (err warehouseError) Error() string {
	return string(err)
}
