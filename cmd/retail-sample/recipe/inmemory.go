package recipe

import domain "github.com/anatollupacescu/retail-sample/internal/retail-domain/recipe"

type InMemoryStore struct {
	data    map[int]domain.Recipe
	counter *int
}

func NewInMemoryStore() InMemoryStore {
	zero := 0
	return InMemoryStore{
		data:    make(map[int]domain.Recipe),
		counter: &zero,
	}
}

func (m *InMemoryStore) Add(r domain.Recipe) (domain.ID, error) {
	*m.counter++
	id := *m.counter
	m.data[id] = r
	return domain.ID(id), nil
}

func (m *InMemoryStore) List() (r []domain.Recipe) {
	for id, rp := range m.data {
		rp.ID = domain.ID(id)
		r = append(r, rp)
	}

	return
}

func (m *InMemoryStore) Get(id domain.ID) domain.Recipe {
	return m.data[int(id)]
}
