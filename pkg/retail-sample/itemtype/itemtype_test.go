package itemtype_test

import (
	"testing"

	"github.com/anatollupacescu/retail-sample/pkg/retail-sample/itemtype"
	qt "github.com/frankban/quicktest"
)

func TestItemTypeRepository(t *testing.T) {

	t.Run("new item type repository is empty", func(t *testing.T) {
		c := qt.New(t)
		d := itemtype.NewInMemoryRepository()
		types := d.List()
		c.Assert(types, qt.HasLen, 0)
	})

	t.Run("can add item", func(t *testing.T) {
		c := qt.New(t)
		d := itemtype.NewInMemoryRepository()
		d.Add("beans")
		types := d.List()
		c.Assert(types, qt.HasLen, 1)
		addedType := types[0]
		c.Assert(addedType, qt.DeepEquals, "beans")
	})

	t.Run("can not add repeated item types", func(t *testing.T) {
		c := qt.New(t)
		d := itemtype.NewInMemoryRepository()
		d.Add("beans")
		d.Add("beans")
		types := d.List()
		c.Assert(types, qt.HasLen, 1)
		addedType := types[0]
		c.Assert(addedType, qt.DeepEquals, "beans")
	})

	t.Run("can remove item types", func(t *testing.T) {
		c := qt.New(t)
		d := itemtype.NewInMemoryRepository()
		id := d.Add("beans")
		d.Remove(id)
		types := d.List()
		c.Assert(types, qt.HasLen, 0)
	})

	t.Run("can get item type by identifier", func(t *testing.T) {
		c := qt.New(t)
		d := itemtype.NewInMemoryRepository()
		id := d.Add("beans")
		c.Assert(id, qt.Equals, uint64(1))
		tp := d.Get(1)
		c.Assert(tp, qt.DeepEquals, "beans")
	})

	t.Run("get non existent type returns zero value", func(t *testing.T) {
		c := qt.New(t)
		d := itemtype.NewInMemoryRepository()
		tp := d.Get(1)
		c.Assert(tp, qt.DeepEquals, "")
	})
}
