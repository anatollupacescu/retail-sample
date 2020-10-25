package stock

import (
	"context"

	"github.com/anatollupacescu/retail-sample/domain/retail/stock"
)

type logger interface {
	Error(string, string, error)
	Info(string, string)
}

func New(ctx context.Context, stock stock.Stock, log logger) Stock {
	return Stock{
		ctx:    ctx,
		stock:  stock,
		logger: log,
	}
}

type Stock struct {
	logger logger
	stock  stock.Stock
	ctx    context.Context
}
