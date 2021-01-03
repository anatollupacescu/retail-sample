package stock

import (
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Save(dto PositionDTO) error {
	return m.Called(dto).Error(0)
}

func (m *MockDB) Get(id int) (PositionDTO, error) {
	args := m.Called(id)
	return args.Get(0).(PositionDTO), args.Error(1)
}
