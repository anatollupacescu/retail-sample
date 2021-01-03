package recipe

import (
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Add(r DTO) (int, error) {
	args := m.Called(r)
	return args.Int(0), args.Error(1)
}

func (m *MockDB) Get(id int) (DTO, error) {
	args := m.Called(id)
	return args.Get(0).(DTO), args.Error(1)
}

func (m *MockDB) Find(n string) (DTO, error) {
	args := m.Called(n)
	return args.Get(0).(DTO), args.Error(1)
}

func (m *MockDB) Save(r DTO) error {
	return m.Called(r).Error(0)
}

type MockInventory struct {
	mock.Mock
}

func (m *MockInventory) Validate(ids ...int) error {
	return m.Called(ids).Error(0)
}
