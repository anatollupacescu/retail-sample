package stock

import (
	"github.com/stretchr/testify/mock"

	"github.com/anatollupacescu/retail-sample/domain/retail/recipe"
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

type MockRecipeDB struct {
	mock.Mock
}

func (m *MockRecipeDB) Get(id int) (recipe.DTO, error) {
	args := m.Called(id)
	return args.Get(0).(recipe.DTO), args.Error(1)
}

type MockValidator struct {
	mock.Mock
}

func (m *MockValidator) Valid(id int) (bool, error) {
	args := m.Called(id)
	return args.Bool(0), args.Error(1)
}
