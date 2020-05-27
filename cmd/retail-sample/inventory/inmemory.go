package inventory

import domain "github.com/anatollupacescu/retail-sample/internal/retail-domain/inventory"

type Entry struct {
	Name    string
	Enabled bool
}

type InMemoryStore struct {
	data    map[int]Entry
	counter *int
}

func NewInMemoryStore() InMemoryStore {
	zero := 0
	return InMemoryStore{
		data:    make(map[int]Entry),
		counter: &zero,
	}
}

func (m *InMemoryStore) Add(s string) (int, error) {
	*m.counter += 1

	newID := *m.counter

	m.data[newID] = Entry{
		Name:    s,
		Enabled: true,
	}

	return newID, nil
}

func (m *InMemoryStore) Find(s string) (int, error) {
	for id, v := range m.data {
		if v.Name == s {
			return id, nil
		}
	}

	return 0, domain.ErrItemNotFound
}

func (m *InMemoryStore) Get(wantedID int) (domain.Item, error) {
	var zeroValueItem domain.Item

	for id, v := range m.data {
		if wantedID == id {
			return domain.Item{
				ID:      id,
				Name:    v.Name,
				Enabled: v.Enabled,
			}, nil
		}
	}

	return zeroValueItem, nil
}

func (m *InMemoryStore) List() (t []domain.Item, err error) {
	for k, v := range m.data {
		t = append(t, domain.Item{
			ID:      k,
			Name:    v.Name,
			Enabled: v.Enabled,
		})
	}

	return
}

func (m *InMemoryStore) Update(i domain.Item) (err error) {
	m.data[i.ID] = Entry{
		Name:    i.Name,
		Enabled: i.Enabled,
	}

	return nil
}
