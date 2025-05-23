package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/alekLukanen/go-templ-htmx-example-app/configure"
)

func main() {

	ctx := context.Background()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	logger.Info("Starting API Service")

	// configure basic Service
	serviceConfiguration, err := configure.NewServiceConfiguration(ctx, logger)
	if err != nil {
		logger.Error("failed to configure service", slog.Any("err", err))
		panic(err)
	}

	serverConfiguration := configure.NewUIHttpConfiguration(ctx, logger, serviceConfiguration)
	logger.Info("configured server", slog.String("address", serverConfiguration.Server.Addr))
	go serverConfiguration.StartServer()

	logger.Info("API Service Running")
	logger.Info("Press Ctrl+C to exit")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	cancelCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	logger.Info("server shutting down")
	if err := serverConfiguration.Close(cancelCtx); err != nil {
		logger.Error("server Shutdown Failed", slog.Any("err", err))
		panic(err)
	}
	logger.Info("server shutdown successfully")

	if err := serviceConfiguration.Close(cancelCtx); err != nil {
		logger.Error("failed to close configured resources")
		panic(err)
	}

	logger.Info("service completely closed")

}
