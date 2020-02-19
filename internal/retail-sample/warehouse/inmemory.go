package warehouse

import "time"

type InMemoryInboundLog map[time.Time]ProvisionEntry

func (i InMemoryInboundLog) Add(k time.Time, v ProvisionEntry) {
	i[k] = v
}

func (i InMemoryInboundLog) List() (r []ProvisionEntry) {
	for _, v := range i {
		r = append(r, v)
	}
	return
}

type InMemoryOutboundConfiguration map[string]OutboundItem

func (m InMemoryOutboundConfiguration) add(o OutboundItem) {
	m[o.Name] = o
}

func (m InMemoryOutboundConfiguration) list() (o []OutboundItem) {
	for _, v := range m {
		o = append(o, v)
	}
	return
}

func (m InMemoryOutboundConfiguration) hasConfig(s string) bool {
	_, f := m[s]
	return f
}

func (m InMemoryOutboundConfiguration) components(s string) []OutboundItemComponent {
	if !m.hasConfig(s) {
		return nil
	}

	return m[s].Items
}

type InMemoryOutboundLog map[time.Time]SoldItem

func (m InMemoryOutboundLog) Add(i SoldItem) {
	m[i.Date] = i
}

func (m InMemoryOutboundLog) List() (r []SoldItem) {
	for _, v := range m {
		r = append(r, v)
	}
	return
}
