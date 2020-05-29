package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ardanlabs/conf"
	kitlog "github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/inventory"
	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/order"
	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/recipe"
	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/stock"
	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/types"

	"github.com/anatollupacescu/retail-sample/internal/version"
)

type serverConfig struct {
	DatabaseURL string `conf:"default:postgres://docker:docker@localhost:5432/retail?pool_max_conns=10"`
	Port        string `conf:"default:8080"`
	DiagPort    string `conf:"default:8081"`
	InMemory    bool   `conf:"default:false"`
}

func main() {
	var config serverConfig

	if err := conf.Parse(os.Args[1:], "RETAIL", &config); err != nil {
		log.Fatalf("parse server configuration values: %v", err)
	}

	businessRouter := mux.NewRouter()

	server := http.Server{
		Addr:    net.JoinHostPort("", config.Port),
		Handler: businessRouter,
	}

	baseLogger := kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr))
	baseLogger = kitlog.With(baseLogger, "ts", kitlog.DefaultTimestampUTC)

	//app

	loggerFactory := newLoggerFactory(baseLogger)

	var persistenceFactory types.PersistenceProviderFactory

	if config.InMemory {
		persistenceFactory = newInMemoryPersistentFactory()
	} else {
		persistenceFactory = newPersistenceFactory(config.DatabaseURL)
	}

	routerLogger := wrapLogger(baseLogger)

	inventory.ConfigureRoutes(businessRouter, routerLogger, loggerFactory, persistenceFactory)
	order.ConfigureRoutes(businessRouter, routerLogger, loggerFactory, persistenceFactory)
	recipe.ConfigureRoutes(businessRouter, routerLogger, loggerFactory, persistenceFactory)
	stock.ConfigureRoutes(businessRouter, routerLogger, loggerFactory, persistenceFactory)

	//static
	businessRouter.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./web/dist"))))

	diagRouter := mux.NewRouter()
	diagRouter.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		if err := persistenceFactory.Ping(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	corsDiagRouter := cors.Default().Handler(diagRouter)

	diag := http.Server{
		Addr:    net.JoinHostPort("", config.DiagPort),
		Handler: corsDiagRouter,
	}

	shutdown := make(chan error, 2)

	logger := kitlog.With(baseLogger,
		"version", version.Version,
		"build_time", version.BuildTime,
		"commit", version.Commit)

	go func() {
		_ = logger.Log("msg", "business logic server is starting", "port", config.Port)

		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			shutdown <- err
		}
	}()

	go func() {
		_ = logger.Log("msg", "diagnostics server is starting", "port", config.DiagPort)

		err := diag.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			shutdown <- err
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case x := <-interrupt:
		_ = logger.Log("msg", "received", "signal", x)

	case err := <-shutdown:
		_ = logger.Log("msg", "received shutdown request", "signal", err)
	}

	timeout, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	err := diag.Shutdown(timeout)
	if err != nil {
		_ = logger.Log("msg", "diagnostic server shutdown failed", "error", err)
	}

	err = server.Shutdown(timeout)
	if err != nil {
		_ = logger.Log("msg", "business server shutdown failed", "error", err)
	}
}
