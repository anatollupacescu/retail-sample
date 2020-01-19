package warehouse

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type MockInventory struct {
	mock.Mock
}

func (m *MockInventory) setQty(s string, number int) {
	_ = m.Called(s, number)
}

func (m *MockInventory) qty(s string) int {
	return m.Called(s).Int(0)
}

func (m *MockInventory) addType(s string) {
	m.Called(s)
}

func (m *MockInventory) hasType(s string) bool {
	args := m.Called(s)
	return args.Bool(0)
}

func (m *MockInventory) types() []string {
	return nil
}

func (m *MockInventory) disable(s string) {
	m.Called(s)
}

//outbound

type MockOutboundConfig struct {
	mock.Mock
}

func (m *MockOutboundConfig) add(i OutboundItem) {
	_ = m.Called(i)
}

func (m *MockOutboundConfig) list() []OutboundItem {
	return nil
}

func (m *MockOutboundConfig) hasConfig(s string) bool {
	return m.Called(s).Bool(0)
}

func (m *MockOutboundConfig) components(s string) []OutboundItemComponent {
	args := m.Called(s)
	res := args.Get(0)
	return res.([]OutboundItemComponent)
}

//inbound log

type MockInboundLog struct {
	mock.Mock
}

func (m *MockInboundLog) Add(t time.Time, i Item) {
	_ = m.Called(t, i)
}

func (m *MockInboundLog) List() []Item {
	return nil
}

//outbound log

type MockOutboundLog struct {
	mock.Mock
}

func (m *MockOutboundLog) Add(i SoldItem) {
	_ = m.Called(i)
}

func (m *MockOutboundLog) List() []SoldItem {
	return nil
}
