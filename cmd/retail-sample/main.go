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

	inmemory "github.com/anatollupacescu/retail-sample/cmd/retail-sample/persistence"
	"github.com/anatollupacescu/retail-sample/domain/version"
	persistence "github.com/anatollupacescu/retail-sample/persistence/postgres"
)

// Configuration is exported to be accesible by the library.
type Configuration struct {
	DatabaseURL string `conf:"default:postgres://docker:docker@localhost:5432/retail?pool_max_conns=10"`
	Port        string `conf:"default:8080"`
	Offline     bool   `conf:"default:false"`
}

func configureRouter(baseLogger kitlog.Logger, config Configuration) *mux.Router {
	routerLogger := middleware.WrapLogger(baseLogger)
	loggerFactory := middleware.BuildNewLoggerFunc(baseLogger)

	/* business */
	persistenceFactory := getPersistenceFactory(config)

	router := mux.NewRouter()

	inventory.ConfigureRoutes(router, routerLogger, loggerFactory, persistenceFactory)
	order.ConfigureRoutes(router, routerLogger, loggerFactory, persistenceFactory)
	recipe.ConfigureRoutes(router, routerLogger, loggerFactory, persistenceFactory)
	stock.ConfigureRoutes(router, routerLogger, loggerFactory, persistenceFactory)

	router.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		if err := persistenceFactory.Ping(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	return router
}

func main() {
	config := getConfig()

	baseLogger := kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stdout))
	baseLogger = kitlog.With(baseLogger, "ts", kitlog.DefaultTimestampUTC)

	router := configureRouter(baseLogger, config)

	corsRouter := cors.AllowAll().Handler(router)

	server := http.Server{
		Addr:    net.JoinHostPort("", config.Port),
		Handler: corsRouter,
	}

	logger := kitlog.With(baseLogger,
		"version", version.Version,
		"build_time", version.BuildTime,
		"commit", version.Commit)

	log := func(args ...interface{}) {
		_ = logger.Log(args)
	}

	log("offline", config.Offline)

	const serverCount = 2
	shutdown := make(chan error, serverCount)

	go func() {
		log("msg", "starting", "port", server.Addr)

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
		log("msg", "received", "signal", x)

	case err := <-shutdown:
		log("msg", "received shutdown request", "error", err)
	}

	const waitForShutdown = 5 * time.Second

	timeout, cancelFunc := context.WithTimeout(context.Background(), waitForShutdown)
	defer cancelFunc()

	if err := server.Shutdown(timeout); err != nil {
		log("msg", "business server shutdown failed", "error", err)
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
		return inmemory.New()
	}

	dbURL := strings.TrimSpace(config.DatabaseURL)

	return persistence.NewPersistenceFactory(dbURL)
}
