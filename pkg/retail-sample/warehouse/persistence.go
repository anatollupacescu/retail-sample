package warehouse

import "errors"

type InMemoryItemDB map[uint64]int

func NewInMemoryDB() InMemoryItemDB {
	return make(map[uint64]int)
}

func (s InMemoryItemDB) Add(id uint64, qty int) {
	s[id] = qty
}

func (s InMemoryItemDB) Update(i uint64, qty int) error {
	if got, ok := s[i]; ok {
		s[i] = got + qty
		return nil
	}
	return errors.New("not found")
}

func (s InMemoryItemDB) Get(i uint64) (int, error) {
	if got, ok := s[i]; ok {
		return got, nil
	}
	return 0, errors.New("not found")
}
