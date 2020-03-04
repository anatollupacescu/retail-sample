package order

import "github.com/stretchr/testify/mock"

//outbound log

type MockOrderStore struct {
	mock.Mock
}

func (m *MockOrderStore) add(i OrderEntry) ID {
	return m.Called(i).Get(0).(ID)
}

func (m *MockOrderStore) all() []OrderEntry {
	return nil
}
