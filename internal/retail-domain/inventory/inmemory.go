package inventory

type Entry struct {
	Name Name
}

type InMemoryStore struct {
	data    map[ID]Entry
	counter *int
}

func NewInMemoryStore() InMemoryStore {
	zero := 0
	return InMemoryStore{
		data:    make(map[ID]Entry),
		counter: &zero,
	}
}

func (m *InMemoryStore) Add(s Name) (ID, error) {
	*m.counter += 1

	newID := ID(*m.counter)

	m.data[newID] = Entry{
		Name: s,
	}

	return newID, nil
}

func (m *InMemoryStore) Find(s Name) (ID, error) {
	for id, v := range m.data {
		if v.Name == s {
			return id, nil
		}
	}

	return ID(0), ErrStoreItemNotFound
}

var zeroValueItem = Item{}

func (m *InMemoryStore) Get(wantedID ID) (Item, error) {
	for id, v := range m.data {
		if wantedID == id {
			return Item{
				ID:   id,
				Name: v.Name,
			}, nil
		}
	}

	return zeroValueItem, nil
}

func (m *InMemoryStore) All() (t []Item, err error) {
	for k, v := range m.data {
		t = append(t, Item{
			ID:   k,
			Name: v.Name,
		})
	}

	return
}
