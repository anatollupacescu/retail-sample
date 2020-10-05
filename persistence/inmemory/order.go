package inmemory

import domain "github.com/anatollupacescu/retail-sample/domain/retail-sample/order"

type Order struct {
	data    map[domain.ID]domain.Order
	counter *int
}

func NewOrder() domain.Store {
	zero := 0

	return Order{
		data:    make(map[domain.ID]domain.Order),
		counter: &zero,
	}
}

func (m Order) Add(i domain.Order) (domain.ID, error) {
	*m.counter++

	newID := domain.ID(*m.counter)
	m.data[newID] = i

	return newID, nil
}

func (m Order) List() (r []domain.Order, err error) {
	for id := range m.data {
		v := m.data[id]
		v.ID = id
		r = append(r, v)
	}

	return
}

func (m Order) Get(id domain.ID) (o domain.Order, err error) {
	return m.data[id], nil
}
