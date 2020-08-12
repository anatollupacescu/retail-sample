package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	kitlog "github.com/go-kit/kit/log"
	"github.com/rs/cors"

	"github.com/ardanlabs/conf"
	"github.com/gorilla/mux"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/app/inventory"
	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/app/order"
	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/app/recipe"
	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/app/stock"
	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/middleware"
	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/persistence"

	"github.com/anatollupacescu/retail-sample/internal/version"
)

type serverConfig struct {
	DatabaseURL string `conf:"default:postgres://docker:docker@localhost:5432/retail?pool_max_conns=10"`
	Port        string `conf:"default:8080"`
	DiagPort    string `conf:"default:8081"`
	InMemory    bool   `conf:"default:false"`
}

func main() {
	baseLogger := kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr))
	baseLogger = kitlog.With(baseLogger, "ts", kitlog.DefaultTimestampUTC)
	routerLogger := middleware.WrapLogger(baseLogger)
	loggerFactory := middleware.BuildNewLoggerFunc(baseLogger)

	config := getConfig()
	persistenceFactory := getPersistenceFactory(config)

	businessRouter := mux.NewRouter()

	inventory.ConfigureRoutes(businessRouter, routerLogger, loggerFactory, persistenceFactory)
	order.ConfigureRoutes(businessRouter, routerLogger, loggerFactory, persistenceFactory)
	recipe.ConfigureRoutes(businessRouter, routerLogger, loggerFactory, persistenceFactory)
	stock.ConfigureRoutes(businessRouter, routerLogger, loggerFactory, persistenceFactory)

	// static site
	businessRouter.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./web/dist"))))

	server := http.Server{
		Addr:    net.JoinHostPort("", config.Port),
		Handler: businessRouter,
	}

	diagRouter := mux.NewRouter()

	diagRouter.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		if err := persistenceFactory.Ping(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	// needs CORS because it runs on a different port
	corsDiagRouter := cors.Default().Handler(diagRouter)

	diag := http.Server{
		Addr:    net.JoinHostPort("", config.DiagPort),
		Handler: corsDiagRouter,
	}

	startServers(&server, &diag, baseLogger)
}

func startServers(server, diag *http.Server, baseLogger kitlog.Logger) {
	logger := kitlog.With(baseLogger,
		"version", version.Version,
		"build_time", version.BuildTime,
		"commit", version.Commit)

	const serverCount = 2
	shutdown := make(chan error, serverCount)

	start := func(name string, server *http.Server) {
		_ = logger.Log("server", name, "msg", "starting", "port", server.Addr)

		if err := server.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				shutdown <- err
			}
		}
	}

	go start("business logic", server)

	go start("diagnostic", diag)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case x := <-interrupt:
		_ = logger.Log("msg", "received", "signal", x)

	case err := <-shutdown:
		_ = logger.Log("msg", "received shutdown request", "error", err)
	}

	const waitForShutdown = 5 * time.Second

	timeout, cancelFunc := context.WithTimeout(context.Background(), waitForShutdown)
	defer cancelFunc()

	if err := diag.Shutdown(timeout); err != nil {
		_ = logger.Log("msg", "diagnostic server shutdown failed", "error", err)
	}

	if err := server.Shutdown(timeout); err != nil {
		_ = logger.Log("msg", "business server shutdown failed", "error", err)
	}
}

func getConfig() (config serverConfig) {
	if err := conf.Parse(os.Args[1:], "RETAIL", &config); err != nil {
		log.Fatalf("parse server configuration values: %v", err)
	}

	return
}

func getPersistenceFactory(config serverConfig) middleware.PersistenceProviderFactory {
	if !config.InMemory {
		return persistence.NewInMemory()
	}

	dbURL := strings.TrimSpace(config.DatabaseURL)

	return persistence.NewPersistenceFactory(dbURL)
}
