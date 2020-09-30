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

// Configuration is exported to be accesible by the library.
type Configuration struct {
	DatabaseURL string `conf:"default:postgres://docker:docker@localhost:5432/retail?pool_max_conns=10"`
	Port        string `conf:"default:8080"`
	Offline     bool   `conf:"default:false"`
}

func main() {
	config := getConfig()

	baseLogger := kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stdout))
	baseLogger = kitlog.With(baseLogger, "ts", kitlog.DefaultTimestampUTC)
	routerLogger := middleware.WrapLogger(baseLogger)
	loggerFactory := middleware.BuildNewLoggerFunc(baseLogger)

	/* business */
	persistenceFactory := getPersistenceFactory(config)

	router := mux.NewRouter()

	inventory.ConfigureRoutes(router, routerLogger, loggerFactory, persistenceFactory)
	order.ConfigureRoutes(router, routerLogger, loggerFactory, persistenceFactory)
	recipe.ConfigureRoutes(router, routerLogger, loggerFactory, persistenceFactory)
	stock.ConfigureRoutes(router, routerLogger, loggerFactory, persistenceFactory)

	/* TODO finish up roles

	if !config.InMemory {
		businessRouter.Use(middleware.Authenticated)
	}
	*/

	router.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		if err := persistenceFactory.Ping(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	corsRouter := cors.AllowAll().Handler(router)

	server := http.Server{
		Addr:    net.JoinHostPort("", config.Port),
		Handler: corsRouter,
	}

	logger := kitlog.With(baseLogger,
		"version", version.Version,
		"build_time", version.BuildTime,
		"commit", version.Commit)

	logger.Log("offline", config.Offline)

	const serverCount = 2
	shutdown := make(chan error, serverCount)

	go func() {
		_ = logger.Log("msg", "starting", "port", server.Addr)

		if err := server.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				shutdown <- err
			}
		}
	}()

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

	if err := server.Shutdown(timeout); err != nil {
		_ = logger.Log("msg", "business server shutdown failed", "error", err)
	}
}

func getConfig() (config Configuration) {
	if err := conf.Parse(os.Args[1:], "RETAIL", &config); err != nil {
		log.Fatalf("parse server configuration values: %v", err)
	}

	return
}

func getPersistenceFactory(config Configuration) middleware.PersistenceProviderFactory {
	if config.Offline {
		return persistence.NewInMemory()
	}

	dbURL := strings.TrimSpace(config.DatabaseURL)

	return persistence.NewPersistenceFactory(dbURL)
}
