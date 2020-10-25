package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"

	"github.com/ardanlabs/conf"
	"github.com/gorilla/mux"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/middleware"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/route/inventory"
	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/route/order"
	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/route/recipe"
	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/route/stock"

	persistence "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/persistence/postgres"
	"github.com/anatollupacescu/retail-sample/internal/version"
)

func main() {
	config := getConfig()

	run(config)
}

// Configuration is exported to be accessible by the library.
type Configuration struct {
	DatabaseURL string `conf:"default:postgres://docker:docker@localhost:5432/retail?pool_max_conns=10"`
	Port        string `conf:"default:8080"`
	Offline     bool   `conf:"default:false"`
}

// nolint: funlen // as this is the main function it can be as long as needed
func run(config Configuration) {
	zLog := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("role", "retail-sample").
		Str("host", os.Getenv("HOST")).
		Logger()

	router := mux.NewRouter()

	router.Use(hlog.NewHandler(zLog),
		hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
			hlog.FromRequest(r).Info().
				Str("method", r.Method).
				Int("status", status).
				Int("size", size).
				Dur("duration", duration).
				Msg("")
		}),
		hlog.RemoteAddrHandler("ip"),
		hlog.UserAgentHandler("user_agent"),
		hlog.RefererHandler("referer"),
		hlog.RequestIDHandler("req_id", "Request-Id"))

	ctx, cancelDB := context.WithCancel(context.Background())
	pool := persistence.NewPersistenceFactory(ctx, config.DatabaseURL)

	router.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		if err := persistence.Ping(pool); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	router.Use(middleware.Transaction(pool), middleware.Recovery)

	router.HandleFunc("/inventory", inventory.GetAll).Methods(http.MethodGet)
	router.HandleFunc("/inventory/{itemID}", inventory.Get).Methods(http.MethodGet)
	router.HandleFunc("/inventory/{itemID}", inventory.Update).Methods(http.MethodPatch)
	router.HandleFunc("/inventory", inventory.Create).Methods(http.MethodPost)

	router.HandleFunc("/order", order.Create).Methods(http.MethodPost)
	router.HandleFunc("/order/{orderID}", order.Get).Methods(http.MethodGet)
	router.HandleFunc("/order", order.GetAll).Methods(http.MethodGet)

	router.HandleFunc("/recipe", recipe.GetAll).Methods(http.MethodGet)
	router.HandleFunc("/recipe/{recipeID}", recipe.Get).Methods(http.MethodGet)
	router.HandleFunc("/recipe/{recipeID}", recipe.Update).Methods(http.MethodPatch)
	router.HandleFunc("/recipe", recipe.Create).Methods(http.MethodPost)

	router.HandleFunc("/stock/provisionlog", stock.GetProvisionLog).Methods(http.MethodGet)

	router.HandleFunc("/stock", stock.GetAll).Methods(http.MethodGet)
	router.HandleFunc("/stock/{itemID}", stock.Get).Methods(http.MethodGet)
	router.HandleFunc("/stock/{itemID}", stock.Provision).Methods(http.MethodPost)

	corsRouter := cors.AllowAll().Handler(router)

	server := http.Server{
		Addr:    net.JoinHostPort("", config.Port),
		Handler: corsRouter,
	}

	bootLogger := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("version", version.Version).
		Str("build_time", version.BuildTime).
		Str("commit", version.Commit).
		Logger()

	bootLogger.Info().Bool("offline", config.Offline)

	const serverCount = 2
	shutdown := make(chan error, serverCount)

	go func() {
		bootLogger.Info().Str("port", server.Addr).Msg("serving")

		if err := server.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				shutdown <- err
			}
		}

		pool.Close()
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case x := <-interrupt:
		bootLogger.Info().Str("signal", x.String()).Msg("received")

	case err := <-shutdown:
		bootLogger.Info().Err(err).Msg("received shutdown request")
	}

	const waitForShutdown = 5 * time.Second

	timeout, cancelServer := context.WithTimeout(context.Background(), waitForShutdown)

	defer cancelServer()

	if err := server.Shutdown(timeout); err != nil {
		bootLogger.Error().Err(err).Msg("business server shutdown failed")
	}

	defer cancelDB()
}

func getConfig() (config Configuration) {
	if err := conf.Parse(os.Args[1:], "RETAIL", &config); err != nil {
		log.Fatalf("parse server configuration values: %v", err)
	}

	return
}
