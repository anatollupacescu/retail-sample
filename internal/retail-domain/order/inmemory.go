package order

type InMemoryStore struct {
	data    map[int]OrderEntry
	counter *int
}

func NewInMemoryStore() store {
	zero := 0
	return InMemoryStore{
		data:    make(map[int]OrderEntry),
		counter: &zero,
	}
}

func (m InMemoryStore) add(i OrderEntry) ID {
	currentID := *m.counter
	m.data[currentID] = i
	*m.counter += 1

	return ID(currentID)
}

func (m InMemoryStore) all() (r []OrderEntry) {
	for _, v := range m.data {
		r = append(r, v)
	}

	return
}
