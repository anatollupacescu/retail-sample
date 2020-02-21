package recipe

type InMemoryStore struct {
	data    map[int]Recipe
	counter *int
}

func NewInMemoryStore() InMemoryStore {
	zero := 0
	return InMemoryStore{
		data:    make(map[int]Recipe),
		counter: &zero,
	}
}

func (m *InMemoryStore) add(r Recipe) (ID, error) {
	*m.counter++
	id := *m.counter
	m.data[id] = r
	return ID(id), nil
}

func (m *InMemoryStore) all() []Recipe {
	return nil
}

func (m *InMemoryStore) get(id ID) Recipe {
	return m.data[int(id)]
}
