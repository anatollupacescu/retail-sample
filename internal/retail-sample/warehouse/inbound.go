package warehouse

import "time"

type (
	Item struct {
		Type string
		Qty  int
	}

	Log interface {
		Add(time.Time, Item)
		List() []Item
	}
)
