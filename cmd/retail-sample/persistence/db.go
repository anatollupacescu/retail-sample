package persistence

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type PgxDB interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}
