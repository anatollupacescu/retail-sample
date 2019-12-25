package warehouse_test

import (
	"github.com/anatollupacescu/retail-sample/internal/retail-sample/warehouse"
)

type TestWrapper struct {
	stock warehouse.Stock
}

func stock() TestWrapper {
	return TestWrapper{
		stock: warehouse.NewStock(),
	}
}

func (t TestWrapper) with(s string, q int) TestWrapper {
	if err := t.stock.ConfigureInboundType(s); err != nil {
		panic(err)
	}
	item := warehouse.Item{Type: s, Qty: q}
	if _, err := t.stock.PlaceInbound(item); err != nil {
		panic(err)
	}
	return t
}

func (t TestWrapper) build() warehouse.Stock {
	return t.stock
}
