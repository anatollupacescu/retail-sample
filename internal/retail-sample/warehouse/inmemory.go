package warehouse

import "time"

type InMemoryInboundLog map[time.Time]Item

func (i InMemoryInboundLog) Add(k time.Time, v Item) {
	i[k] = v
}

func (i InMemoryInboundLog) List() (r []Item) {
	for _, v := range i {
		r = append(r, v)
	}
	return
}

type InMemoryInventory map[string]int

func (m InMemoryInventory) setQty(s string, i int) {
	m[s] = i
}

func (m InMemoryInventory) qty(s string) int {
	return m[s]
}

func (m InMemoryInventory) addType(s string) {
	m[s] = 0
}

func (m InMemoryInventory) hasType(s string) bool {
	_, f := m[s]
	return f
}

func (m InMemoryInventory) types() (t []string) {
	for k := range m {
		t = append(t, k)
	}

	return
}

type InMemoryOutboundConfiguration map[string]OutboundItem

func (m InMemoryOutboundConfiguration) add(o OutboundItem) {
	m[o.name] = o
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
