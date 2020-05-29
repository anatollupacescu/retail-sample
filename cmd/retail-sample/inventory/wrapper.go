package inventory

import (
	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/types"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/inventory"
)

type wrapper struct {
	loggerFactory              types.LoggerFactory
	persistenceProviderFactory types.PersistenceProviderFactory
}

func (ia wrapper) exec(methodName string, f func(i inventory.Inventory) error) {
	logger := ia.loggerFactory()

	logger.Log("msg", "enter", "method", methodName)
	defer logger.Log("msg", "exit", "method", methodName)

	provider := ia.persistenceProviderFactory.New()
	inventory := provider.Inventory()

	err := f(inventory)

	if err != nil {
		logger.Log("msg", "rollback")
		ia.persistenceProviderFactory.Rollback(provider)
		return
	}

	logger.Log("msg", "commit")
	ia.persistenceProviderFactory.Commit(provider)
}

func (ia wrapper) setStatus(id int, enabled bool) (item inventory.Item, err error) {
	ia.exec("update item status", func(i inventory.Inventory) error {
		item, err = i.UpdateStatus(id, enabled)

		return err
	})

	return
}

func (ia wrapper) create(name string) (id int, err error) {
	ia.exec("add to inventory", func(i inventory.Inventory) error {
		itemName := name

		id, err = i.Add(itemName)

		return err
	})

	return
}

func (ia wrapper) getAll() (items []inventory.Item, err error) {
	ia.exec("list inventory items", func(i inventory.Inventory) error {
		items, err = i.List()

		return err
	})

	return
}

func (ia wrapper) getOne(id int) (item inventory.Item, err error) {
	ia.exec("get inventory item", func(i inventory.Inventory) error {
		item, err = i.Get(id)

		return err
	})

	return
}
