package acceptance_test

import (
	"errors"
	"fmt"

	"github.com/anatollupacescu/arbortest/runner"

	http "github.com/anatollupacescu/retail-sample/cmd/retail-sample-test"
	random "github.com/anatollupacescu/retail-sample/cmd/retail-sample-test"

	domain "github.com/anatollupacescu/retail-sample/domain/retail-sample/stock"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/app/stock"
)

func provisionStock(id, qty int) (int, error) {
	req := http.PostToID("stock", id)

	return stock.Provision(qty, req)
}

func getStockPosition(id int) (domain.Position, error) {
	gcl := http.Get("stock", id)

	return stock.Get(gcl)
}

// group:stock after:inventory
func testProvision(t *runner.T) {
	createdName := random.Name()
	item, _ := createItem(createdName)

	itemID := item.ID

	var (
		reqQty = 9
		newQty int
	)

	newQty, err := provisionStock(itemID, reqQty)

	if err != nil {
		t.Error(err)
		return
	}

	if newQty != reqQty {
		t.Error(errors.New("should have the provisioned quantity"))
	}
}

// group:stock
func testGetStockPos(t *runner.T) {
	createdName := random.Name()
	item, _ := createItem(createdName)

	itemID := item.ID

	var reqQty = 9

	_, _ = provisionStock(itemID, reqQty)

	pos, err := getStockPosition(itemID)

	if err != nil {
		t.Error(err)
		return
	}

	if pos.Qty != reqQty {
		t.Error(fmt.Errorf("expected qty 0, got %v", pos.Qty))
	}
}

// group:stock
func testGetAllStockPos(t *runner.T) {
	createdName := random.Name()
	item, _ := createItem(createdName)

	itemID := item.ID

	var reqQty = 9

	_, _ = provisionStock(itemID, reqQty)

	gcl := http.List("stock")

	all, err := stock.GetAll(gcl)

	if err != nil {
		t.Error(err)
		return
	}

	var found domain.Position

	for _, v := range all {
		if v.ID == itemID {
			found = v
		}
	}

	if found.ID == 0 {
		t.Error(errors.New("item not found in stock"))
		return
	}

	if found.Qty != reqQty {
		t.Error(fmt.Errorf("expected quantity %d, got %d", reqQty, found.Qty))
	}
}
