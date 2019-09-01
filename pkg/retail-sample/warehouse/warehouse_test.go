package warehouse_test

import (
	"testing"

	"github.com/anatollupacescu/retail-sample/pkg/retail-sample/itemtype"
	"github.com/anatollupacescu/retail-sample/pkg/retail-sample/warehouse"
	qt "github.com/frankban/quicktest"
)

func NewTestingWarehouse() warehouse.Repository {
	return warehouse.Repository{
		ItemDB:       warehouse.NewInMemoryDB(),
		ItemTypeRepo: itemtype.NewInMemoryRepository(),
	}
}

func TestWarehouse(t *testing.T) {

	t.Run("should reject non existent item type ids", func(t *testing.T) {
		c := qt.New(t)
		w := NewTestingWarehouse()
		err := w.Add(1, 23)
		c.Assert(err, qt.Equals, warehouse.ErrItemTypeNotFound)
	})

	t.Run("should add item if item type is present", func(t *testing.T) {
		c := qt.New(t)
		wr := NewTestingWarehouse()
		id := wr.ItemTypeRepo.Add("test")

		err := wr.Add(id, 23)
		c.Assert(err, qt.IsNil)

		qty, err := wr.Quantity(id)
		c.Assert(err, qt.IsNil)
		c.Assert(qty, qt.Equals, 23)
	})

	t.Run("should return error when querying quantity of non existent item type", func(t *testing.T) {
		c := qt.New(t)
		wr := NewTestingWarehouse()
		_, err := wr.Quantity(uint64(1))
		c.Assert(err, qt.Equals, warehouse.ErrItemTypeNotFound)
	})

	t.Run("should add different item types", func(t *testing.T) {
		c := qt.New(t)
		wr := NewTestingWarehouse()
		itr := wr.ItemTypeRepo

		idTest := itr.Add("test")
		idSecond := itr.Add("second item type")

		wr.Add(idTest, 23)
		wr.Add(idSecond, 49)

		qty, err := wr.Quantity(idTest)
		c.Assert(err, qt.IsNil)
		c.Assert(qty, qt.Equals, 23)

		qty, err = wr.Quantity(idSecond)
		c.Assert(err, qt.IsNil)
		c.Assert(qty, qt.Equals, 49)
	})

	t.Run("should combine quantities for the same type", func(t *testing.T) {
		c := qt.New(t)
		wr := NewTestingWarehouse()
		itr := wr.ItemTypeRepo

		typeID := itr.Add("test")

		wr.Add(typeID, 23)
		wr.Add(typeID, 49)

		qty, err := wr.Quantity(typeID)
		c.Assert(err, qt.IsNil)
		c.Assert(qty, qt.Equals, 23+49)
	})
}
