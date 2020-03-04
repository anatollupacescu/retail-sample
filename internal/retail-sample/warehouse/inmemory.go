package warehouse

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

type InMemoryOrderLog map[time.Time]OrderLogEntry

func (m InMemoryOrderLog) Add(i OrderLogEntry) {
	m[i.Date] = i
}

func (m InMemoryOrderLog) List() (r []OrderLogEntry) {
	for _, v := range m {
		r = append(r, v)
	}
	return
}
