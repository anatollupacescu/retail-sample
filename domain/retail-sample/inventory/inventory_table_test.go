package inventory_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/anatollupacescu/retail-sample/domain/retail-sample/inventory"
)

func TestAddT(t *testing.T) {
	var tt = []struct {
		testName string
		store    testStore
		name     string
		//expected
		id  int
		err error
	}{
		{
			"rejects empty name",
			testStore{},
			"",
			0, inventory.ErrEmptyName,
		}, {
			"rejects duplicate name",
			testStore{
				find: func(s string) (int, error) {
					return 1, nil
				},
			},
			"test",
			0, inventory.ErrDuplicateName,
		}, {
			"when call to store's 'add' throws error, it's propagated",
			testStore{
				find: func(s string) (int, error) {
					return 0, inventory.ErrItemNotFound
				},
				add: func(s string) (int, error) {
					return 0, errors.New("unexpected")
				},
			},
			"test",
			0, errors.New("unexpected"),
		}, {
			"adds a valid name",
			testStore{
				find: func(s string) (int, error) {
					return 0, inventory.ErrItemNotFound
				},
				add: func(s string) (int, error) {
					assert.Equal(t, "test", s)
					return 5, nil
				},
			},
			"test",
			5, nil,
		},
	}

	for i := 0; i < len(tt); i++ {
		test := tt[i]

		i := inventory.Inventory{Store: &test.store}

		id, err := i.Add(test.name)

		assert.Equal(t, test.err, err, test.testName)
		assert.Equal(t, test.id, id, test.testName)
	}
}

func TestUpdateStatusT(t *testing.T) {
	var tt = []struct {
		testName string
		store    testStore
		//inputs
		id      int
		enabled bool
		//expected
		item inventory.Item
		err  error
		//calls
		calls string
	}{
		{
			"when item is not found in store it propagates the error",
			testStore{
				get: func(int) (inventory.Item, error) {
					return inventory.Item{}, inventory.ErrItemNotFound
				},
			},
			1, true,
			inventory.Item{}, inventory.ErrItemNotFound, "get",
		}, {
			"when store call to 'Get' throws unknown error it is propagated",
			testStore{
				get: func(int) (inventory.Item, error) {
					return inventory.Item{}, errors.New("unexpected")
				},
			},
			1, true,
			inventory.Item{}, errors.New("unexpected"), "get",
		}, {
			"when store call to 'Update' throws unknown error it is propagated",
			testStore{
				get: func(int) (inventory.Item, error) {
					return inventory.Item{ID: 99}, nil
				},
				update: func(i inventory.Item) error {
					return errors.New("unexpected")
				},
			},
			1, true,
			inventory.Item{}, errors.New("unexpected"), "getupdate",
		}, {
			"when item is found its status is changed",
			testStore{
				get: func(int) (inventory.Item, error) {
					return inventory.Item{Enabled: false}, nil
				},
				update: func(i inventory.Item) error {
					assert.Equal(t, inventory.Item{Enabled: true}, i)
					return nil
				},
			},
			1, true,
			inventory.Item{Enabled: true}, nil, "getupdate",
		},
	}

	for i := 0; i < len(tt); i++ {
		test := tt[i]

		i := inventory.Inventory{Store: &test.store}

		item, err := i.UpdateStatus(test.id, test.enabled)

		assert.Equal(t, test.err, err, test.testName)
		assert.Equal(t, test.item, item, test.testName)
		assert.Equal(t, test.store.calls, test.calls, test.testName)
	}
}

type testStore struct {
	calls  string
	add    func(string) (int, error)
	find   func(string) (int, error)
	get    func(int) (inventory.Item, error)
	list   func() ([]inventory.Item, error)
	update func(inventory.Item) error
}

func (t *testStore) Add(s string) (int, error) {
	t.calls += "add"
	return t.add(s)
}

func (t *testStore) Find(s string) (int, error) {
	t.calls += "find"
	return t.find(s)
}

func (t *testStore) Get(i int) (inventory.Item, error) {
	t.calls += "get"
	return t.get(i)
}

func (t *testStore) List() ([]inventory.Item, error) {
	t.calls += "list"
	return t.list()
}

func (t *testStore) Update(i inventory.Item) error {
	t.calls += "update"
	return t.update(i)
}
