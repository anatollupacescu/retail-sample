// +build acceptance

package arbor_test

import (
	"errors"
	"fmt"
	"testing"

	invCmd "github.com/anatollupacescu/retail-sample/cmd/retail-sample/app/inventory"
	"github.com/anatollupacescu/retail-sample/internal/arbor"

	"github.com/stretchr/testify/assert"

	"github.com/anatollupacescu/retail-sample/internal/retail-domain/inventory"

	"github.com/google/uuid"
)

var baseURL = "http://localhost:8080/inventory"

func TestInventory(t *testing.T) {
	add := arbor.New("create item", func() error {
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
	})

	duplicatesRejected := arbor.New("name duplicate rejected", func() error {
		name := uuid.New().String()

		_, _ = invCmd.Create(baseURL, name)

		if _, err := invCmd.Create(baseURL, name); err == nil {
			return errors.New("expected error")
		}

		return nil
	}, add)

	duplicatesRejected.Run()

	disable := arbor.New("disable item", func() error {
		name := uuid.New().String()

		_, _ = invCmd.Create(baseURL, name)

		if _, err := invCmd.Create(baseURL, name); err == nil {
			return errors.New("expected error")
		}

		return nil
	}, add)

	disable.Run()

	fmt.Print(disable.String())

	t.Run("succeeds", func(t *testing.T) {
		assert.Equal(t, "[create item] ok\n", add.String())
		assert.Equal(t, "[name duplicate rejected] ok\n⮡ [create item] ok\n", duplicatesRejected.String())
		assert.Equal(t, "[disable] ok\n⮡ [create item] ok\n", disable.String())
	})
}
