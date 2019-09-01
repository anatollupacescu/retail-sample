package warehouse

import "github.com/anatollupacescu/retail-sample/pkg/retail-sample/itemtype"

type (
	itemTypeID uint64
	Repository map[itemTypeID]int
)

const (
	ErrItemTypeNotFound = warehouseError("Item type with given id was not found")
	ErrItemNotFound = warehouseError("No such item stored in the warehouse")
)

var (
	ItemTypeRepository = itemtype.NewRepository()
	zeroItemTypeValue = itemtype.ItemType{}
)

func (r *Repository) Add(id uint64, qty int) error {
	itemType := ItemTypeRepository.Get(id)
	if itemType == zeroItemTypeValue {
		return ErrItemTypeNotFound
	}
	if got, ok := (*r)[itemTypeID(id)]; ok {
		(*r)[itemTypeID(id)] = got + qty
		return nil
	}
	(*r)[itemTypeID(id)] = qty
	return nil
}

func (r *Repository) Quantity(i uint64) (int ,error) {
	wantedItemType := ItemTypeRepository.Get(i)
	if wantedItemType == zeroItemTypeValue {
		return 0, ErrItemTypeNotFound
	}
	if qty, ok := (*r)[itemTypeID(i)]; ok {
		return qty, nil
	}
	return 0, ErrItemNotFound
}

type warehouseError string

func (err warehouseError) Error() string {
	return string(err)
}
