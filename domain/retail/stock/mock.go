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
