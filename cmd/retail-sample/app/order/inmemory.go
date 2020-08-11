package order

import domain "github.com/anatollupacescu/retail-sample/internal/retail-domain/order"

type InMemoryStore struct {
	data    map[domain.ID]domain.Order
	counter *int
}

func NewInMemoryStore() domain.Store {
	zero := 0

	return InMemoryStore{
		data:    make(map[domain.ID]domain.Order),
		counter: &zero,
	}
}

func (m InMemoryStore) Add(i domain.Order) (domain.ID, error) {
	*m.counter++

	newID := domain.ID(*m.counter)
	m.data[newID] = i

	return newID, nil
}

func (m InMemoryStore) List() (r []domain.Order, err error) {
	for id := range m.data {
		v := m.data[id]
		v.ID = id
		r = append(r, v)
	}

	return
}

func (m InMemoryStore) Get(id domain.ID) (o domain.Order, err error) {
	return m.data[id], nil
}
