package middleware

import (
	"context"
	"net/http"

	persistence "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/persistence/postgres"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

// nolint:gocognit // transaction handing middleware is bound to have this complexity
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

			// add to context to make it available for downstream
			ctx := context.WithValue(rCtx, TxKey, tx)
			r = r.WithContext(ctx)

			crw := &capturingResponseWriter{w, http.StatusOK}

			// rollback transaction in case of panic
			defer recoverTransaction(ctx, crw, tx, logger)

			next.ServeHTTP(crw, r)

			if crw.hasErrored() {
				if err := tx.Rollback(ctx); err != nil {
					logger.Err(err).Msg("rollback transaction")
					return
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

func (w *capturingResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

const httpErrorTier = 400

func (w *capturingResponseWriter) hasErrored() bool {
	return w.statusCode >= httpErrorTier
}

func recoverTransaction(ctx context.Context, w http.ResponseWriter, tx persistence.TX, logger *zerolog.Logger) {
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
