package arbor_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	invCmd "github.com/anatollupacescu/retail-sample/cmd/retail-sample/app/inventory"
	"github.com/anatollupacescu/retail-sample/internal/arbor"

	"github.com/anatollupacescu/retail-sample/internal/retail-domain/inventory"

	"github.com/google/uuid"
)

var baseURL = "http://localhost:8080/inventory"

func testCreate() error {
	var item inventory.Item
	var err error

	name := uuid.New().String()

	if item, err = invCmd.Create(baseURL, name); err != nil {
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
	name := uuid.New().String()

	_, _ = invCmd.Create(baseURL, name)

	if _, err := invCmd.Create(baseURL, name); err == nil {
		return errors.New("expected error")
	}

	return nil
}

func testDisable() error {
	name := uuid.New().String()

	_, _ = invCmd.Create(baseURL, name)

	if _, err := invCmd.Create(baseURL, name); err == nil {
		return errors.New("expected error")
	}

	return nil
}

func TestInventory(t *testing.T) {
	create := arbor.New("create item", testCreate)
	// getOne
	// getAll
	noDuplicate := arbor.New("name duplicate rejected", testDuplicate, create)
	disable := arbor.New("disable item", testDisable, create)

	all := arbor.New("all", func() error { return nil }, noDuplicate, disable)

	all.Run()

	t.Run("succeeds", func(t *testing.T) {
		assert.Equal(t, true, all.Success)
	})
}
