package acceptance_test

import (
	"errors"
	"fmt"

	http "github.com/anatollupacescu/retail-sample/cmd/retail-sample-test"

	domain "github.com/anatollupacescu/retail-sample/internal/retail-domain/stock"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/app/stock"
)

func provisionStock(id, qty int) (int, error) {
	req := http.PostToID("stock", id)

	return stock.Provision(qty, req)
}

func testProvision() (err error) {
	itemID := createRandomItem()

	var (
		reqQty = 9
		newQty int
	)

	newQty, err = provisionStock(itemID, reqQty)

	if err != nil {
		return err
	}

	if newQty != reqQty {
		return errors.New("should have the provisioned quantity")
	}

	return nil
}

func getStockPos(id int) (domain.Position, error) {
	gcl := http.Get("stock", id)

	return stock.Get(gcl)
}

func testGetStockPos() error {
	itemID := createRandomItem()

	var reqQty = 9

	_, _ = provisionStock(itemID, reqQty)

	pos, err := getStockPos(itemID)

	if err != nil {
		return err
	}

	if pos.Qty != reqQty {
		return fmt.Errorf("expected qty 0, got %v", pos.Qty)
	}

	return nil
}

func testGetAllStockPos() error {
	itemID := createRandomItem()

	var reqQty = 9

	_, _ = provisionStock(itemID, reqQty)

	gcl := http.List("stock")

	all, err := stock.GetAll(gcl)

	if err != nil {
		return err
	}

	var found domain.Position

	for _, v := range all {
		if v.ID == itemID {
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
