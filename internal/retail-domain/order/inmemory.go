package order

type InMemoryStore struct {
	data    map[ID]Order
	counter *int
}

func NewInMemoryStore() store {
	zero := 0
	return InMemoryStore{
		data:    make(map[ID]Order),
		counter: &zero,
	}
}

func (m InMemoryStore) add(i Order) ID {
	*m.counter += 1

	newID := ID(*m.counter)
	m.data[newID] = i

	return newID
}

func (m InMemoryStore) all() (r []Order) {
	for id, v := range m.data {
		v.ID = id
		r = append(r, v)
	}

	return
}
