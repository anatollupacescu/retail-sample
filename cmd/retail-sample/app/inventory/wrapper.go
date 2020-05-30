package inventory

import (
	types "github.com/anatollupacescu/retail-sample/cmd/retail-sample/middleware"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/inventory"
)

type wrapper struct {
	types.Wrapper
}

func (ia wrapper) setStatus(id int, enabled bool) (item inventory.Item, err error) {
	return item, ia.Exec("update item status", func(provider types.PersistenceProvider) error {
		i := provider.Inventory()

		item, err = i.UpdateStatus(id, enabled)

		return err
	})
}

func (ia wrapper) create(name string) (id int, err error) {
	ia.Exec("add to inventory", func(provider types.PersistenceProvider) error {
		i := provider.Inventory()

		id, err = i.Add(name)

		return err
	})

	return
}

func (ia wrapper) getAll() (items []inventory.Item, err error) {
	ia.Exec("list inventory items", func(provider types.PersistenceProvider) error {
		i := provider.Inventory()

		items, err = i.List()

		return err
	})

	return
}

func (ia wrapper) getOne(id int) (item inventory.Item, err error) {
	ia.Exec("get inventory item", func(provider types.PersistenceProvider) error {
		i := provider.Inventory()

		item, err = i.Get(id)

		return err
	})

	return
}
