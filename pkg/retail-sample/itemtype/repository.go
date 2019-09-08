package itemtype

func NewInMemoryRepository() *InMemoryDB {
	return &InMemoryDB{
		data:    make(map[string]uint64),
		counter: 0,
	}
}
