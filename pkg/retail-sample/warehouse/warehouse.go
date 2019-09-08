//go:generate mockgen -source=warehouse.go -package mocks -destination mocks/warehouse.go
package warehouse

type (
	ItemRepository interface {
		Add(uint64, int)
		Update(uint64, int) error
		Get(uint64) (int, error)
	}

	ItemTypeRepository interface {
		Add(string) uint64
		Get(uint64) string
		Remove(uint64)
		Find(string) (uint64, error)
		List() []string
	}

	Repository struct {
		ItemRepository     ItemRepository
		ItemTypeRepository ItemTypeRepository
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

func (r *Repository) Get(id uint64) (int, error) {
	return r.ItemRepository.Get(id)
}

func (r *Repository) Add(id uint64, qty int) error {
	itemType := r.ItemTypeRepository.Get(id)

	if itemType == zeroItemTypeValue {
		return ErrItemTypeNotFound
	}

	if _, err := r.ItemRepository.Get(id); err != nil {
		r.ItemRepository.Add(id, qty)
		return nil
	}

	if err := r.ItemRepository.Update(id, qty); err != nil {
		return ErrUpdate
	}

	return nil
}

type warehouseError string

func (err warehouseError) Error() string {
	return string(err)
}
