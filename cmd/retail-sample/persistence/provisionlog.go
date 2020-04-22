package persistence

import retail "github.com/anatollupacescu/retail-sample/internal/retail-sample"

type PgxProvisionLog struct {
	DB PgxDB
}

func (po *PgxProvisionLog) Add(retail.ProvisionEntry) {
}

func (po *PgxProvisionLog) List() []retail.ProvisionEntry {
	return nil
}
