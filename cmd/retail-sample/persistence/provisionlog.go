package persistence

import (
	"context"

	retail "github.com/anatollupacescu/retail-sample/internal/retail-sample"
	"github.com/pkg/errors"
)

type PgxProvisionLog struct {
	DB PgxDB
}

func (po *PgxProvisionLog) Add(re retail.ProvisionEntry) error {
	sql := "insert into provisionlog(inventoryid, quantity) values($1, $2)"
	if _, err := po.DB.Exec(context.Background(), sql, re.ID, re.Qty); err != nil {
		return errors.Wrapf(DBErr, "provisionlog add: %v", err)
	}

	return nil
}

func (po *PgxProvisionLog) List() (ee []retail.ProvisionEntry, err error) {
	rows, err := po.DB.Query(context.Background(), "select inventoryid, quantity from provisionlog")

	if err != nil {
		return nil, errors.Wrapf(DBErr, "provisionlog list: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var id int64
		var qty int16

		if err := rows.Scan(&id, &qty); err != nil {
			return nil, errors.Wrapf(DBErr, "provisionlog list scan: %v", err)
		}

		ee = append(ee, retail.ProvisionEntry{
			ID:  int(id),
			Qty: int(qty),
		})
	}

	return
}
