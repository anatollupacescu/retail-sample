package middleware

import (
	"context"
	"errors"
	"net/http"

	persistence "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/persistence/postgres"
)

type Tx int

const TxKey = Tx(0)

var (
	errTransactionNotFound = errors.New("transaction not found")
	errTransactionBadType  = errors.New("transaction invalid")
)

func ExtractTransaction(r *http.Request) (tx persistence.TX, err error) {
	ctxTransaction := r.Context().Value(TxKey)

	if ctxTransaction == nil {
		return persistence.TX{}, errTransactionNotFound
	}

	var ok bool

	if tx, ok = ctxTransaction.(persistence.TX); !ok {
		return persistence.TX{}, errTransactionBadType
	}

	return
}

func ExtractTransactionCtx(ctx context.Context) (tx persistence.TX, err error) {
	ctxTransaction := ctx.Value(TxKey)

	if ctxTransaction == nil {
		return persistence.TX{}, errTransactionNotFound
	}

	var ok bool

	if tx, ok = ctxTransaction.(persistence.TX); !ok {
		return persistence.TX{}, errTransactionBadType
	}

	return
}
