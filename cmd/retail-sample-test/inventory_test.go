package acceptance_test

import (
	"errors"

	http "github.com/anatollupacescu/retail-sample/cmd/retail-sample-test"
	random "github.com/anatollupacescu/retail-sample/cmd/retail-sample-test"
	web "github.com/anatollupacescu/retail-sample/cmd/retail-sample/app/inventory"
	domain "github.com/anatollupacescu/retail-sample/internal/retail-domain/inventory"
)

func createRandomItem() int {
	name := random.Name()
	i, _ := createItem(name)

	return i.ID
}

func createItem(name string) (domain.Item, error) {
	cl := http.Post("inventory")

	return web.Create(name, cl)
}

func testCreate() (err error) {
	name := random.Name()

	item, err := createItem(name)

	if err != nil {
		return err
	}

	if item.Name != name {
		return errors.New("bad name")
	}

	if item.ID == 0 {
		return errors.New("bad id")
	}

	return nil
}

func testCreateWithEmptyName() error {
	_, err := createItem("")

	if err == nil {
		return errors.New("should return nil for empty name")
	}

	return nil
}

func testDuplicate() error {
	name := random.Name()

	_, _ = createItem(name)

	_, err := createItem(name)

	if err == nil {
		return errors.New("should reject duplicate name")
	}

	return nil
}

func testDisable() error {
	itemID := createRandomItem()

	cl := http.Patch("inventory", itemID)

	updatedItem, err := web.Update(false, cl)

	if err != nil {
		return err
	}

	if updatedItem.Enabled != false {
		return errors.New("should disable item")
	}

	return nil
}

func testGetAll() (err error) { //TODO create an item an assert it's present in the 'all'
	cl := http.List("inventory")

	all, err := web.GetAll(cl)

	if err != nil {
		return err
	}

	if len(all) < 1 {
		return errors.New("should return multiple items")
	}

	return nil
}

func testGetOne() (err error) {
	name := random.Name()

	i, _ := createItem(name)

	gcl := http.Get("inventory", i.ID)

	item, err := web.Get(gcl)

	if err != nil {
		return err
	}

	if item.Name != name {
		return errors.New("should have the same name")
	}

	if item.ID == 0 {
		return errors.New("should not have zero value for ID")
	}

	return nil
}
