package retailsampleapp1

import (
	"time"
)

type InMemoryProvisionLog map[time.Time]ProvisionEntry

func (i InMemoryProvisionLog) Add(v ProvisionEntry) {
	i[time.Now()] = v
}

func (i InMemoryProvisionLog) List() (r []ProvisionEntry) {
	for t, v := range i {
		r = append(r, ProvisionEntry{
			Time: t,
			ID:   v.ID,
			Qty:  v.Qty,
		})
	}
	return
}
