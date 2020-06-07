// +build acceptance

package inv_test

import (
	"errors"
	"flag"
	"fmt"
	"testing"
	"time"

	faker "github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"

	"github.com/anatollupacescu/retail-sample/internal/arbor"

	client "github.com/anatollupacescu/retail-sample/cmd/retail-sample-test"
	web "github.com/anatollupacescu/retail-sample/cmd/retail-sample/app/inventory"
	domain "github.com/anatollupacescu/retail-sample/internal/retail-domain/inventory"
)

var (
	graphURL = flag.String("arborURL", "", "graph server URL")
	apiURL   = flag.String("apiURL", "", "api server URL")

	timeout = 100 * time.Millisecond //TODO pass as flag
)

func TestInventory(t *testing.T) {
	if *apiURL == "" {
		t.Fatal("api URL not provided")
	}

	create := arbor.New("create", testCreate)
	getOne := arbor.New("get one", testGetOne, create)
	getAll := arbor.New("get all", testGetAll, create)
	noDuplicate := arbor.New("no duplicate", testDuplicate, create)
	disable := arbor.New("disable", testDisable, create)

	all := arbor.Suite("all", getAll, getOne, noDuplicate, disable)

	all.Run()

	t.Run("succeeds", func(t *testing.T) {
		assert.Equal(t, true, all.Success)
	})

	t.Logf("\n%s", all)

	js := arbor.Marshal(create, getOne, getAll, noDuplicate, disable)

	arbor.Upload(*graphURL, js)
}

func testCreate() (err error) {
	var item domain.Item

	name := faker.Word()

	cl := client.Post(*apiURL, timeout)

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

	cl := client.Post(*apiURL, timeout)
	_, _ = web.Create(name, cl)

	if _, err := web.Create(name, cl); err == nil {
		return errors.New("expected error")
	}

	return nil
}

func testDisable() (err error) {
	name := faker.Word()

	cl := client.Post(*apiURL, timeout)
	i, _ := web.Create(name, cl)

	resourceURL := fmt.Sprintf("%s/%d", *apiURL, i.ID)

	cl = client.Patch(resourceURL, timeout)

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
	cl := client.Get(*apiURL, timeout)

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

	cl := client.Post(*apiURL, timeout)
	i, _ := web.Create(name, cl)

	resourceURL := fmt.Sprintf("%s/%d", *apiURL, i.ID)
	gcl := client.Get(resourceURL, timeout)

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
