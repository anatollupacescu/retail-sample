package warehouse

import "github.com/stretchr/testify/mock"

type MockedInventory struct{
	mock.Mock
}

func (m MockedInventory) setQty(s string, number int) {
	_ = m.Called(s, number)
}

func (m MockedInventory) qty(string) int {
	return 0
}

func (m MockedInventory) addType(string) {

}

func (m MockedInventory) hasType(s string) bool {
	args := m.Called(s)
	return args.Bool(0)
}

func (m MockedInventory) types() []string {
	return nil
}

