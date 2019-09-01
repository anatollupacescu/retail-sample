package itemtype_test

import (
	"github.com/anatollupacescu/retail-sample/pkg/retail-sample/itemtype"
	qt "github.com/frankban/quicktest"
	"testing"
)

func TestItemTypeRepository(t *testing.T) {

	t.Run("new item type repository is empty", func(t *testing.T) {
		c := qt.New(t)
		d := itemtype.Repository{}
		types := d.ListItemTypes()
		c.Assert(types, qt.HasLen, 0)
	})

	t.Run("can add item", func(t *testing.T) {
		c := qt.New(t)
		d := itemtype.Repository{}
		d.AddItemType("beans", 3)
		types := d.ListItemTypes()
		c.Assert(types, qt.HasLen, 1)
		addedType := types[0]
		c.Assert(addedType, qt.DeepEquals, itemtype.ItemType{Name: "beans", Qty: 3})
	})

	t.Run("can not add repeated item types", func(t *testing.T) {
		c := qt.New(t)
		d := itemtype.Repository{}
		d.AddItemType("beans", 3)
		d.AddItemType("beans", 3)
		types := d.ListItemTypes()
		c.Assert(types, qt.HasLen, 1)
		addedType := types[0]
		c.Assert(addedType, qt.DeepEquals, itemtype.ItemType{Name: "beans", Qty: 3})
	})

	t.Run("can remove item types", func(t *testing.T) {
		c := qt.New(t)
		d := itemtype.Repository{}
		d.AddItemType("beans", 3)
		d.RemoveItemType("beans", 3)
		types := d.ListItemTypes()
		c.Assert(types, qt.HasLen, 0)
	})
}
