package inventory

import (
	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/middleware"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/inventory"
)

type wrapper struct {
	middleware.Wrapper
}

func (ia wrapper) setStatus(id int, enabled bool) (item inventory.Item, err error) {
	return item, ia.Exec("update item status", func(provider middleware.PersistenceProvider) error {
		i := provider.Inventory()

		item, err = i.UpdateStatus(id, enabled)

		return err
	})
}

func (ia wrapper) create(name string) (item inventory.Item, err error) {
	return item, ia.Exec("add to inventory", func(provider middleware.PersistenceProvider) error {
		i := provider.Inventory()

		var id int

		id, err = i.Add(name)

		if err != nil {
			return err
		}

		item, err = i.Get(id)

		return err
	})
}

func (ia wrapper) getAll() (items []inventory.Item, err error) {
	return items, ia.Exec("list inventory items", func(provider middleware.PersistenceProvider) error {
		i := provider.Inventory()

		items, err = i.List()

		return err
	})
}

func (ia wrapper) getOne(id int) (item inventory.Item, err error) {
	return item, ia.Exec("get inventory item", func(provider middleware.PersistenceProvider) error {
		i := provider.Inventory()

		item, err = i.Get(id)

		return err
	})
}
