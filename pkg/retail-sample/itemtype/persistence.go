package itemtype

import (
	"errors"
	"sync/atomic"
)

type InMemoryDB struct {
	data    map[string]uint64
	counter uint64
}

func (db *InMemoryDB) Add(name string) uint64 {
	id := atomic.AddUint64(&db.counter, 1)
	db.data[name] = id
	return id
}

func (db *InMemoryDB) Get(i uint64) string {
	for itemType, gotID := range db.data {
		if i == gotID {
			return itemType
		}
	}
	return ""
}

func (db *InMemoryDB) Remove(i uint64) {
	t := db.Get(i)
	if t != "" {
		delete(db.data, t)
	}
}

func (db *InMemoryDB) Find(t string) (i uint64, err error) {
	if _, ok := db.data[t]; !ok {
		return 0, errors.New("not found")
	}
	return db.data[t], nil
}

func (db *InMemoryDB) List() []string {
	types := make([]string, 0, len(db.data))
	for t := range db.data {
		types = append(types, t)
	}
	return types
}
