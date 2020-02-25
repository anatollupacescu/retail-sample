package warehouse

import (
	"time"
)

type InMemoryInboundLog map[time.Time]ProvisionEntry

func (i InMemoryInboundLog) Add(v ProvisionEntry) {
	i[time.Now()] = v
}

func (i InMemoryInboundLog) List() (r []ProvisionEntry) {
	for t, v := range i {
		r = append(r, ProvisionEntry{
			Time: t,
			ID:   v.ID,
			Qty:  v.Qty,
		})
	}
	return
}

type InMemoryOutboundLog map[time.Time]OrderLogEntry

func (m InMemoryOutboundLog) Add(i OrderLogEntry) {
	m[i.Date] = i
}

func (m InMemoryOutboundLog) List() (r []OrderLogEntry) {
	for _, v := range m {
		r = append(r, v)
	}
	return
}
