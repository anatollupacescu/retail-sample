package itemtype_test

import (
	"github.com/anatollupacescu/retail-sample/pkg/retail-sample/itemtype"
	qt "github.com/frankban/quicktest"
	"testing"
)

func TestItemTypeRepository(t *testing.T) {

	t.Run("new item type repository is empty", func(t *testing.T) {
		c := qt.New(t)
		d := itemtype.NewRepository()
		types := d.List()
		c.Assert(types, qt.HasLen, 0)
	})

	t.Run("can add item", func(t *testing.T) {
		c := qt.New(t)
		d := itemtype.NewRepository()
		d.Add("beans")
		types := d.List()
		c.Assert(types, qt.HasLen, 1)
		addedType := types[0]
		c.Assert(addedType, qt.DeepEquals, itemtype.ItemType{Name: "beans"})
	})

	t.Run("can not add repeated item types", func(t *testing.T) {
		c := qt.New(t)
		d := itemtype.NewRepository()
		d.Add("beans")
		d.Add("beans")
		types := d.List()
		c.Assert(types, qt.HasLen, 1)
		addedType := types[0]
		c.Assert(addedType, qt.DeepEquals, itemtype.ItemType{Name: "beans"})
	})

	t.Run("can remove item types", func(t *testing.T) {
		c := qt.New(t)
		d := itemtype.NewRepository()
		d.Add("beans")
		d.RemoveItemType("beans", 3)
		types := d.List()
		c.Assert(types, qt.HasLen, 0)
	})

	t.Run("can get item type by identifier", func(t *testing.T) {
		c := qt.New(t)
		d := itemtype.NewRepository()
		id := d.Add("beans")
		c.Assert(id, qt.Equals, uint64(1))
		tp := d.Get(1)
		c.Assert(tp, qt.DeepEquals, itemtype.ItemType{Name:"beans"})
	})
}
