package persistence

import (
	"context"
	"log"

	retail "github.com/anatollupacescu/retail-sample/internal/retail-sample"
)

type PgxProvisionLog struct {
	DB PgxDB
}

func (po *PgxProvisionLog) Add(re retail.ProvisionEntry) {
	sql := "insert into provisionlog(inventoryid, quantity) values($1, $2)"
	if _, err := po.DB.Exec(context.Background(), sql, re.ID, re.Qty); err != nil {
		log.Print("provisionlog add ", err)
	}
}

func (po *PgxProvisionLog) List() (ee []retail.ProvisionEntry) {
	rows, err := po.DB.Query(context.Background(), "select inventoryid, quantity from provisionlog")

	if err != nil {
		log.Print("provisionlog list ", err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var id int64
		var qty int16
		if err := rows.Scan(&id, &qty); err != nil {
			log.Print("provisionlog list scan ", err)
			return
		}
		ee = append(ee, retail.ProvisionEntry{
			ID:  int(id),
			Qty: int(qty),
		})
	}

	return
}
