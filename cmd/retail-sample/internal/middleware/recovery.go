package middleware

import (
	"errors"
	"net/http"

	persistence "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/persistence/postgres"

	"github.com/rs/zerolog/hlog"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer concludeTransaction(w, r)

		next.ServeHTTP(w, r)
	})
}

func concludeTransaction(w http.ResponseWriter, r *http.Request) {
	logger := hlog.FromRequest(r)

	tx, err := ExtractTransaction(r)

	if err != nil {
		logger.Err(err).Msg("get transaction from context")
		http.Error(w, "server error", http.StatusInternalServerError)

		return
	}

	if rr := recover(); rr != nil {
		logger.Error().Msgf("panic: %v", rr)
		http.Error(w, "server error", http.StatusInternalServerError)

		if err := tx.Rollback(r.Context()); err != nil {
			logger.Err(err).Msg("rollback transaction")

			return
		}

		logger.Info().Msg("rollback successful")

		return
	}

	if err := tx.Commit(r.Context()); err != nil {
		logger.Err(err).Msg("commit transaction")
		http.Error(w, "server error", http.StatusInternalServerError)

		return
	}

	logger.Info().Msg("commit successful")
}

var (
	errTransactionNotFound = errors.New("transaction not found")
	errTransactionBadType  = errors.New("transaction not postgres")
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
