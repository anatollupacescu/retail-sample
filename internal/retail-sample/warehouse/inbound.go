package warehouse

import "time"

type (
	ItemType string

	Item struct {
		Type ItemType
		Qty  int
	}

	Log interface {
		Add(time.Time, Item)
		List() []Item
	}

	InMemoryInboundLog map[time.Time]Item
)

func (i InMemoryInboundLog) Add(k time.Time, v Item) {
	i[k] = v
}

func (i InMemoryInboundLog) List() (r []Item) {
	for _, v := range i {
		r = append(r, v)
	}
	return
}
