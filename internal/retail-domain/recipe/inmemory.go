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

func (m *InMemoryStore) Add(r Recipe) (ID, error) {
	*m.counter++
	id := *m.counter
	m.data[id] = r
	return ID(id), nil
}

func (m *InMemoryStore) List() (r []Recipe) {
	for id, rp := range m.data {
		rp.ID = ID(id)
		r = append(r, rp)
	}

	return
}

func (m *InMemoryStore) Get(id ID) Recipe {
	return m.data[int(id)]
}
