package inventory

import (
	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/types"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/inventory"
)

type InventoryWrapper struct {
	LoggerFactory              types.LoggerFactory
	PersistenceProviderFactory types.PersistenceProviderFactory
}

func (ia InventoryWrapper) exec(methodName string, f func(i inventory.Inventory) error) {
	logger := ia.LoggerFactory()

	logger.Log("msg", "enter", "method", methodName)
	defer logger.Log("msg", "exit", "method", methodName)

	provider := ia.PersistenceProviderFactory.New()
	inventory := provider.Inventory()

	err := f(inventory)

	if err != nil {
		logger.Log("msg", "rollback")
		ia.PersistenceProviderFactory.Rollback(provider)
		return
	}

	logger.Log("msg", "commit")
	ia.PersistenceProviderFactory.Commit(provider)
}

func (ia InventoryWrapper) AddToInventory(name string) (id int, err error) {
	ia.exec("add to inventory", func(i inventory.Inventory) error {
		itemName := name

		id, err = i.Add(itemName)

		return err
	})

	return
}

func (ia InventoryWrapper) ListInventoryItems() (items []inventory.Item, err error) {
	ia.exec("list inventory items", func(i inventory.Inventory) error {
		items, err = i.List()

		return err
	})

	return
}

func (ia InventoryWrapper) GetInventoryItem(id int) (item inventory.Item, err error) {
	ia.exec("get inventory item", func(i inventory.Inventory) error {
		item, err = i.Get(id)

		return err
	})

	return
}
