package itemtype

import "sync/atomic"

type (
	Entity struct {
		name string
	}

	InMemoryDB struct {
		data    map[Entity]uint64
		counter uint64
	}
)

var (
	zeroEntity = Entity{}
)

func (db *InMemoryDB) Add(s string) uint64 {
	id := atomic.AddUint64(&db.counter, 1)
	db.data[Entity{name: s}] = id
	return id
}

func (db *InMemoryDB) Get(i uint64) Entity {
	for itemType, gotID := range db.data {
		if i == gotID {
			return itemType
		}
	}
	return zeroEntity
}

func (db *InMemoryDB) Remove(i uint64) {
	t := db.Get(i)
	if t != zeroEntity {
		delete(db.data, t)
	}
}

func (db *InMemoryDB) List() []Entity {
	types := make([]Entity, 0, len(db.data))
	for t := range db.data {
		types = append(types, t)
	}
	return types
}

func NewInMemoryRepository() Repository {
	return Repository{
		DB: &InMemoryDB{
			data:    make(map[Entity]uint64),
			counter: 0,
		},
	}
}
