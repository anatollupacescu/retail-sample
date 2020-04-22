package persistence

import "github.com/anatollupacescu/retail-sample/internal/retail-domain/order"

type PgxOrderStore struct {
	DB PgxDB
}

func (po *PgxOrderStore) Add(order.Order) order.ID {
	return order.ID(0)
}

func (po *PgxOrderStore) List() []order.Order {
	return nil
}
