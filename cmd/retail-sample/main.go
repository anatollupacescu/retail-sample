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

	"github.com/gorilla/mux"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/web"
	"github.com/anatollupacescu/retail-sample/internal/version"

	"github.com/rs/cors"

	kitlog "github.com/go-kit/kit/log"
)

func newGoKitLogger() kitlog.Logger {
	var logger kitlog.Logger
	logger = kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr))
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)
	return logger
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Business logic port is not set")
	}

	diagPort := os.Getenv("DIAG_PORT")
	if diagPort == "" {
		log.Fatal("Diagnostics port is not set")
	}

	businessRouter := mux.NewRouter()

	server := http.Server{
		Addr:    net.JoinHostPort("", port),
		Handler: businessRouter,
	}

	baseLogger := newGoKitLogger()

	logger := kitlog.With(baseLogger, "version", version.Version,
		"build_time", version.BuildTime,
		"commit", version.Commit)

	logger.Log("msg", "the application is starting")

	appLogger := kitlog.With(baseLogger, "caller", kitlog.DefaultCaller)
	webApp := web.NewApp(appLogger)

	//app
	web.ConfigureRoutes(businessRouter, webApp)

	//static
	businessRouter.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./web/dist"))))

	diagRouter := mux.NewRouter()
	diagRouter.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		var status = http.StatusOK

		if !webApp.IsHealthy() {
			status = http.StatusInternalServerError
		}

		w.WriteHeader(status)
	})

	diagRouter.HandleFunc("/ready", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	corsDiagRouter := cors.Default().Handler(diagRouter)

	diag := http.Server{
		Addr:    net.JoinHostPort("", diagPort),
		Handler: corsDiagRouter,
	}

	shutdown := make(chan error, 2)

	go func() {
		logger.Log("msg", "business logic server is starting", "port", port)

		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			shutdown <- err
		}
	}()

	go func() {
		logger.Log("msg", "diagnostics server is starting", "port", diagPort)

		err := diag.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			shutdown <- err
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case x := <-interrupt:
		logger.Log("msg", "received", "signal", x)

	case err := <-shutdown:
		logger.Log("msg", "received shutdown request", "signal", err)
	}

	timeout, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	err := diag.Shutdown(timeout)
	if err != nil {
		logger.Log("msg", "diagnostic server shutdown failed", "error", err)
	}

	err = server.Shutdown(timeout)
	if err != nil {
		logger.Log("msg", "business server shutdown failed", "error", err)
	}
}
