// +build acceptance

package acceptance_test

import (
	"errors"
	"flag"
	"testing"

	faker "github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"

	"github.com/anatollupacescu/retail-sample/internal/arbor"

	client "github.com/anatollupacescu/retail-sample/cmd/retail-sample-test"
	web "github.com/anatollupacescu/retail-sample/cmd/retail-sample/app/inventory"
	domain "github.com/anatollupacescu/retail-sample/internal/retail-domain/inventory"
)

var arborURL = flag.String("arborURL", "", "graph server URL")

func TestInventory(t *testing.T) {
	createEmpty := arbor.New("rejects empty name", testCreateWithEmptyName)
	createOk := arbor.New("can create", testCreate)

	create := arbor.Suite("create", createEmpty, createOk)

	getOne := arbor.New("get one", testGetOne, create)
	getAll := arbor.New("get all", testGetAll, create)
	noDuplicate := arbor.New("no duplicate", testDuplicate, create)
	disable := arbor.New("disable", testDisable, create)

	all := arbor.Suite("all", getAll, getOne, noDuplicate, disable)

	all.Run()

	t.Run("succeeds", func(t *testing.T) {
		assert.Equal(t, true, all.Success)
	})

	t.Logf("%s\n", all)

	report := arbor.Marshal(create, createEmpty, createOk, getOne, getAll, noDuplicate, disable)

	arbor.Upload(*arborURL, report)
}

func testCreateWithEmptyName() (err error) {
	cl := client.Post()

	if _, err = web.Create("", cl); err == nil {
		return errors.New("expected err")
	}

	return nil
}

func testCreate() (err error) {
	name := faker.Word()

	cl := client.Post()

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
	name := faker.Word()

	cl := client.Post()

	_, _ = web.Create(name, cl)

	if _, err := web.Create(name, cl); err == nil {
		return errors.New("expected error")
	}

	return nil
}

func testDisable() (err error) {
	name := faker.Word()

	cl := client.Post()

	i, _ := web.Create(name, cl)

	cl = client.Patch(i.ID)

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
	cl := client.Get()

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
	name := faker.Word()

	cl := client.Post()

	i, _ := web.Create(name, cl)

	gcl := client.Get(i.ID)

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
