package inventory

type Entry struct {
	Name     Name
	Disabled bool
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

func (m InMemoryStore) add(s Name) ID {
	*m.counter += 1

	newID := ID(*m.counter)

	m.data[newID] = Entry{
		Name: s,
	}

	return newID
}

func (m InMemoryStore) find(s Name) ID {
	for id, v := range m.data {
		if v.Name == s {
			return id
		}
	}

	return ID(0)
}

func (m InMemoryStore) all() (t []Record) {
	for k, v := range m.data {
		t = append(t, Record{
			ID:   k,
			Name: v.Name,
		})
	}

	return
}
