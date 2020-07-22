package acceptance_test

import (
	"errors"
	"fmt"

	"github.com/anatollupacescu/arbortest/runner"

	http "github.com/anatollupacescu/retail-sample/cmd/retail-sample-test"
	random "github.com/anatollupacescu/retail-sample/cmd/retail-sample-test"
	web "github.com/anatollupacescu/retail-sample/cmd/retail-sample/app/inventory"
	domain "github.com/anatollupacescu/retail-sample/internal/retail-domain/inventory"
)

func createItem(name string) (domain.Item, error) {
	cl := http.Post("inventory")

	return web.Create(name, cl)
}

// group:inventory
func testCreate(t *runner.T) {
	name := random.Name()

	item, err := createItem(name)

	if err != nil {
		t.Error(err)
		return
	}

	if item.Name != name {
		t.Error(errors.New("bad name"))
		return
	}

	if item.ID == 0 {
		t.Error(errors.New("bad id"))
	}
}

// group:inventory
func testCreateWithEmptyName(t *runner.T) {
	_, err := createItem("")

	if err == nil {
		t.Error(errors.New("should return error for empty name"))
	}
}

// group:inventory
func testDuplicate(t *runner.T) {
	name := random.Name()
	_, _ = createItem(name)
	_, err := createItem(name)

	if err == nil {
		t.Error(errors.New("should reject duplicate name"))
	}
}

// group:inventory
func testDisable(t *runner.T) {
	name := random.Name()
	item, _ := createItem(name)

	cl := http.Patch("inventory", int(item.ID))

	updatedItem, err := web.Update(false, cl)

	if err != nil {
		t.Error(err)
		return
	}

	if updatedItem.Enabled != false {
		t.Error(errors.New("should disable item"))
	}
}

// group:inventory
func testGetAll(t *runner.T) { //TODO create an item and assert it's present in the 'all'
	name := random.Name()
	item, _ := createItem(name)

	cl := http.List("inventory")

	all, err := web.GetAll(cl)

	if err != nil {
		t.Error(err)
		return
	}

	if len(all) < 1 {
		t.Error(errors.New("should return multiple items"))
		return
	}

	for _, i := range all {
		if i.ID == item.ID {
			return
		}
	}

	t.Error(fmt.Errorf("created item with ID %q wasn't found in response", item.ID))
}

// group:inventory
func testGetOne(t *runner.T) {
	name := random.Name()
	i, _ := createItem(name)

	gcl := http.Get("inventory", i.ID)

	item, err := web.Get(gcl)

	if err != nil {
		t.Error(err)
		return
	}

	if item.Name != name {
		t.Error(errors.New("should have the same name"))
		return
	}

	if item.ID == 0 {
		t.Error(errors.New("should not have zero value for ID"))
	}
}
