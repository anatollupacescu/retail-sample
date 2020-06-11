// +build acceptance

package acceptance_test

import (
	"errors"
	"fmt"

	domain "github.com/anatollupacescu/retail-sample/internal/retail-domain/stock"

	http "github.com/anatollupacescu/retail-sample/cmd/retail-sample-test"
	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/app/inventory"
	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/app/stock"

	faker "github.com/bxcodec/faker/v3"
)

func testProvision() (err error) {
	name := faker.Word()

	req := http.Post("inventory")

	item, _ := inventory.Create(name, req)

	req = http.Post("stock", item.ID)

	var (
		reqQty = 9
		newQty int
	)

	if newQty, err = stock.Provision(reqQty, req); err != nil {
		return err
	}

	if newQty != reqQty {
		return errors.New("bad quantity")
	}

	// can fetch the provisioned item

	gcl := http.Get("stock", item.ID)

	pos, err := stock.Get(gcl)

	if err != nil {
		return err
	}

	if pos.Qty != reqQty {
		return fmt.Errorf("expected qty %d, got %d", reqQty, pos.Qty)
	}

	return nil
}

func testGetStockPos() error {
	name := faker.Word()

	cl := http.Post("inventory")

	item, _ := inventory.Create(name, cl)

	gcl := http.Get("stock", item.ID)

	pos, err := stock.Get(gcl)

	if err != nil {
		return err
	}

	if pos.Qty != 0 {
		return fmt.Errorf("expected qty 0, got %v", pos.Qty)
	}

	if pos.Name != name {
		return fmt.Errorf("expected name '%s', got '%s'", name, pos.Name)
	}

	return nil
}

func testGetAllStockPos() error {
	name := faker.Word()

	cl := http.Post("inventory")

	item, _ := inventory.Create(name, cl)

	gcl := http.Get("stock")

	all, err := stock.GetAll(gcl)

	if err != nil {
		return err
	}

	var foundPos domain.Position

	for _, v := range all {
		if v.ID == item.ID {
			foundPos = v
			break
		}
	}

	if foundPos.ID == 0 {
		return errors.New("item not found in stock")
	}

	return nil
}
