package inventory

import "errors"

type ( //inventory

	Name string
	ID   int

	Item struct {
		ID   ID
		Name Name
	}

	Store interface {
		Add(Name) (ID, error)
		Find(Name) (ID, error)
		Get(ID) (Item, error)
		List() ([]Item, error)
	}

	Inventory struct {
		Store Store
	}
)

var (
	ErrStoreItemNotFound     = errors.New("item not found")
	ErrInventoryItemNotFound = errors.New("item not found")

	ErrEmptyName     = errors.New("name not provided")
	ErrDuplicateName = errors.New("item type already present")

	zeroID = ID(0)
)

func (i Inventory) Add(name Name) (ID, error) {
	if name == "" {
		return zeroID, ErrEmptyName
	}

	_, err := i.Store.Find(name)

	switch err {
	case ErrStoreItemNotFound: //success
		return i.Store.Add(name)
	case nil:
		return zeroID, ErrDuplicateName
	default:
		return zeroID, err
	}
}

func (i Inventory) List() ([]Item, error) {
	return i.Store.List()
}

func (i Inventory) Find(name Name) (ID, error) {
	return i.Store.Find(name)
}

func (i Inventory) Get(id ID) (Item, error) {
	item, err := i.Store.Get(id)

	if err == ErrStoreItemNotFound {
		return zeroValueItem, ErrInventoryItemNotFound
	}

	return item, err
}
