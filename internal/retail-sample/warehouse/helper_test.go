package warehouse_test

import "github.com/anatollupacescu/retail-sample/internal/retail-sample/warehouse"

type TestWrapper struct {
	stock warehouse.Stock
}

func stock() TestWrapper {
	return TestWrapper{
		stock: warehouse.NewStock(),
	}
}

func (t TestWrapper) with(n string, q int) TestWrapper {
	if err := t.stock.ConfigureInboundType(n); err != nil {
		panic(err)
	}
	item := warehouse.InboundItem{Type: warehouse.InboundType(n), Qty: q}
	if _, err := t.stock.PlaceInbound(item); err != nil {
		panic(err)
	}
	return t
}

func (t TestWrapper) build() warehouse.Stock {
	return t.stock
}
