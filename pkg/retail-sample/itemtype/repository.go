package itemtype

import "sync/atomic"

type (
	repository struct {
		data map[ItemType]uint64
		counter uint64
	}

	ItemType struct {
		Name string
	}
)

func NewRepository() repository {
	return repository{
		data:    make(map[ItemType]uint64),
		counter: 0,
	}
}

func (r *repository) List() []ItemType {
	types := make([]ItemType, 0, len(r.data))
	for t := range r.data {
		types = append(types, t)
	}
	return types
}

func (r *repository) Add(name string) uint64 {
	atomic.AddUint64(&r.counter, 1)
	r.data[ItemType{Name: name}] = r.counter
	return r.counter
}

func (r *repository) RemoveItemType(name string, qty int) {
	delete(r.data, ItemType{Name: name})
}

func (r *repository) Get(i uint64) ItemType {
	for itemType, gotID := range r.data {
		if i == gotID {
			return itemType
		}
	}
	return ItemType{}
}
