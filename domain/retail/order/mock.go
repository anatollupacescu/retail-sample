package order

import (
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func (db *MockDB) Add(i DTO) (int, error) {
	args := db.Called(i)
	return args.Int(0), args.Error(1)
}

type MockStock struct {
	mock.Mock
}

func (s *MockStock) Extract(id, qty int) error {
	return s.Called(id, qty).Error(0)
}

type MockRecipe struct {
	mock.Mock
}

func (r *MockRecipe) Valid(id int) error {
	return r.Called(id).Error(0)
}
