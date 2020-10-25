package middleware

import (
	"context"
	"net/http"

	persistence "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/persistence/postgres"
)

type Tx int

const TxKey = Tx(0)

func Transaction(db *persistence.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rCtx := r.Context()
			tx, err := db.Begin(rCtx)

			if err != nil {
				http.Error(w, "could not start transaction", http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(rCtx, TxKey, tx)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
