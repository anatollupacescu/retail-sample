//go:generate mockgen -source=warehouse.go -package mocks -destination mocks/warehouse.go
package warehouse

import "github.com/anatollupacescu/retail-sample/pkg/retail-sample/itemtype"

type (
	itemTypeID uint64

	ItemStore interface {
		Add(uint64, int)
		Update(uint64, int) error
		Get(uint64) (int, error)
	}

	Repository struct {
		ItemStore          ItemStore
		ItemTypeRepository itemtype.Repository
	}
)

const (
	ErrItemTypeNotFound = warehouseError("Item type with given id was not found")
	ErrItemNotFound     = warehouseError("No such item stored in the warehouse")
)

var (
	zeroItemTypeValue = itemtype.ItemType{}
)

func (r *Repository) Add(id uint64, qty int) error {
	itemType := r.ItemTypeRepository.Get(id)
	if itemType == zeroItemTypeValue {
		return ErrItemTypeNotFound
	}
	if _, err := r.ItemStore.Get(id); err == nil {
		r.ItemStore.Update(id, qty)
		return nil
	}
	r.ItemStore.Add(id, qty)
	return nil
}

func (r *Repository) Quantity(i uint64) (int, error) {
	wantedItemType := r.ItemTypeRepository.Get(i)
	if wantedItemType == zeroItemTypeValue {
		return 0, ErrItemTypeNotFound
	}
	if qty, err := r.ItemStore.Get(i); err == nil {
		return qty, nil
	}
	return 0, ErrItemNotFound
}

type warehouseError string

func (err warehouseError) Error() string {
	return string(err)
}
