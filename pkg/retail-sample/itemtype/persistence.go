package itemtype

import "sync/atomic"

type (
	InMemoryDB struct {
		data    map[string]uint64
		counter uint64
	}
)

var (
	zeroDTO = DTO{}
)

func (db *InMemoryDB) Add(name string) DTO {
	id := atomic.AddUint64(&db.counter, 1)
	dto := DTO{Name: name, Id:id}
	db.data[name] = id
	return dto
}

func (db *InMemoryDB) Get(i uint64) DTO {
	for itemType, gotID := range db.data {
		if i == gotID {
			return DTO{Name: itemType, Id:i}
		}
	}
	return zeroDTO
}

func (db *InMemoryDB) Remove(i uint64) {
	t := db.Get(i)
	if t != zeroDTO {
		delete(db.data, t.Name)
	}
}

func (db *InMemoryDB) List() []DTO {
	types := make([]DTO, 0, len(db.data))
	for k,v := range db.data {
		types = append(types, DTO{Name: k, Id:v})
	}
	return types
}

func NewInMemoryRepository() Repository {
	store := &InMemoryDB{
		data:    make(map[string]uint64),
		counter: 0,
	}
	return &repository{
		store: store,
	}
}
