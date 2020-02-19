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

type InMemoryInventory struct {
	config  map[string]bool
	data    map[string]int
	ids     map[string]int
	counter *int
}

func (m InMemoryInventory) setQty(s string, i int) {
	m.data[s] = i
}

func (m InMemoryInventory) qty(s string) int {
	return m.data[s]
}

func (m InMemoryInventory) addType(s string) int {
	m.data[s] = 0
	m.config[s] = false
	*m.counter += 1
	m.ids[s] = *m.counter
	return *m.counter
}

func (m InMemoryInventory) hasType(s string) bool {
	_, f := m.config[s]

	return f
}

func (m InMemoryInventory) types() (t []ItemConfig) {
	for k := range m.data {
		t = append(t, ItemConfig{
			Type:     k,
			Disabled: m.config[k],
		})
	}

	return
}

func (m InMemoryInventory) disable(s string) {
	if !m.hasType(s) {
		return
	}

	m.config[s] = true
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
