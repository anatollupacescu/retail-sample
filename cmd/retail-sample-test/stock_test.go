package acceptance_test

import (
	"errors"
	"fmt"

	http "github.com/anatollupacescu/retail-sample/cmd/retail-sample-test"
	random "github.com/anatollupacescu/retail-sample/cmd/retail-sample-test"

	domain "github.com/anatollupacescu/retail-sample/internal/retail-domain/stock"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/app/inventory"
	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/app/stock"
)

func testProvision() (err error) {
	name := random.Word()

	req := http.Post("inventory")

	item, _ := inventory.Create(name, req)

	req = http.Add("stock", item.ID)

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
	name := random.Word()

	cl := http.Post("inventory")

	item, _ := inventory.Create(name, cl)

	//provision

	req := http.Add("stock", item.ID)

	var reqQty = 9

	_, _ = stock.Provision(reqQty, req)

	// get

	gcl := http.Get("stock", item.ID)

	pos, err := stock.Get(gcl)

	if err != nil {
		return err
	}

	if pos.Qty != reqQty {
		return fmt.Errorf("expected qty 0, got %v", pos.Qty)
	}

	if pos.Name != name {
		return fmt.Errorf("expected name '%s', got '%s'", name, pos.Name)
	}

	return nil
}

func testGetAllStockPos() error {
	name := random.Word()

	cl := http.Post("inventory")

	item, _ := inventory.Create(name, cl)

	//provision

	req := http.Add("stock", item.ID)

	var reqQty = 9

	_, _ = stock.Provision(reqQty, req)

	gcl := http.List("stock")

	all, err := stock.GetAll(gcl)

	if err != nil {
		return err
	}

	var found domain.Position

	for _, v := range all {
		if v.ID == item.ID {
			found = v
		}
	}

	if found.ID == 0 {
		return errors.New("item not found in stock")
	}

	if found.Qty != reqQty {
		return fmt.Errorf("expected quantity %d, got %d", reqQty, found.Qty)
	}

	return nil
}
