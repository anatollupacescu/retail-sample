package middleware

import (
	"context"
	"errors"
	"net/http"

	persistence "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/persistence/postgres"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

type Tx int

const TxKey = Tx(0)

func WithTransaction(db *persistence.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := hlog.FromRequest(r)
			rCtx := r.Context()
			tx, err := db.Begin(rCtx)

			// start a new transaction
			if err != nil {
				logger.Err(err).Msg("could not start transaction")
				http.Error(w, "server error", http.StatusInternalServerError)
				return
			}

			// save it in the context to make it available for downstream
			ctx := context.WithValue(rCtx, TxKey, tx)
			r = r.WithContext(ctx)

			crw := &capturingResponseWriter{w, http.StatusOK}

			// rollback transaction in case of panic
			defer recoverTransaction(crw, tx, ctx, logger)

			next.ServeHTTP(crw, r)

			if crw.statusCode >= 400 {
				if err := tx.Rollback(ctx); err != nil {
					logger.Err(err).Msg("rollback transaction")
				}

				logger.Info().Msg("rollback successful")

				return
			}

			if err := tx.Commit(ctx); err != nil {
				logger.Err(err).Msg("commit transaction")
				http.Error(w, "server error", http.StatusInternalServerError)

				return
			}

			logger.Info().Msg("commit successful")
		})
	}
}

type capturingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *capturingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func recoverTransaction(w http.ResponseWriter, tx persistence.TX, ctx context.Context, logger *zerolog.Logger) {
	if rr := recover(); rr != nil {
		logger.Error().Msgf("panic: %v", rr)
		http.Error(w, "server error", http.StatusInternalServerError)

		if err := tx.Rollback(ctx); err != nil {
			logger.Err(err).Msg("rollback transaction")
			return
		}

		logger.Info().Msg("rollback successful")
	}
}

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
