// +build acceptance

package acceptance_test

import (
	"errors"
	"fmt"

	random "github.com/anatollupacescu/retail-sample/cmd/retail-sample-test"

	http "github.com/anatollupacescu/retail-sample/cmd/retail-sample-test"
	web "github.com/anatollupacescu/retail-sample/cmd/retail-sample/app/inventory"
	domain "github.com/anatollupacescu/retail-sample/internal/retail-domain/inventory"
)

func testCreateWithEmptyName() (err error) {
	cl := http.Post("inventory")

	if _, err = web.Create("", cl); err == nil {
		return errors.New("expected err")
	}

	return nil
}

func testCreate() (err error) {
	name := random.Word()

	cl := http.Post("inventory")

	var item domain.Item

	if item, err = web.Create(name, cl); err != nil {
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

func testDuplicate() error {
	name := random.Word()

	cl := http.Post("inventory")

	_, _ = web.Create(name, cl)

	if _, err := web.Create(name, cl); err == nil {
		return errors.New("expected error")
	}

	return nil
}

func testDisable() (err error) {
	name := random.Word()

	cl := http.Post("inventory")

	i, err := web.Create(name, cl)

	if err != nil {
		return fmt.Errorf("could not create prereq inv item: %v", err)
	}

	cl = http.Patch("inventory", i.ID)

	var updatedItem domain.Item

	if updatedItem, err = web.Update(false, cl); err != nil {
		return err
	}

	if updatedItem.Enabled != false {
		return errors.New("expected resource to be updated")
	}

	return nil
}

func testGetAll() (err error) { //TODO create an item an assert it's present in the 'all'
	cl := http.Get("inventory")

	all, err := web.GetAll(cl)

	if err != nil {
		return err
	}

	if len(all) < 1 {
		return errors.New("expected more items")
	}

	return nil
}

func testGetOne() (err error) {
	name := random.Word()

	cl := http.Post("inventory")

	i, _ := web.Create(name, cl)

	gcl := http.Get("inventory", i.ID)

	item, err := web.Get(gcl)

	if err != nil {
		return err
	}

	if item.Name != name {
		return errors.New("bad name")
	}

	if item.ID == 0 {
		return errors.New("bad ID")
	}

	return nil
}
