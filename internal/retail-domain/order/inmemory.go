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

func (m InMemoryStore) Add(i Order) (ID, error) {
	*m.counter += 1

	newID := ID(*m.counter)
	m.data[newID] = i

	return newID, nil
}

func (m InMemoryStore) List() (r []Order, err error) {
	for id, v := range m.data {
		v.ID = id
		r = append(r, v)
	}

	return
}

func (m InMemoryStore) Get(id ID) (o Order, err error) {
	return m.data[id], nil
}
