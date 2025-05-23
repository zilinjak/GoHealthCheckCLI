package main

import (
	"GoHealthChecker/internal"
	"GoHealthChecker/internal/controller"
	"GoHealthChecker/internal/model"
	"GoHealthChecker/internal/service"
	"GoHealthChecker/internal/store"
	"GoHealthChecker/internal/view"
	"context"
	"os"
	"os/signal"
)

func signalHandler() (context.Context, context.CancelFunc) {
	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	// Set up signal handling
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	// Handle signals in a separate goroutine
	go func() {
		<-signalCh
		cancel() // Cancel context on Ctrl+C
	}()

	return ctx, cancel
}

func main() {
	TIMEOUT := 10
	POOLING_INTERVAL := 1
	ctx, _ := signalHandler()

	settings := model.AppSettings{
		Timeout:         TIMEOUT,
		PoolingInterval: POOLING_INTERVAL,
		Context:         ctx,
	}

	inMemoryStore := store.NewInMemoryStore()
	CLIView := view.NewCLIView()
	HTTPService := service.NewHTTPService(settings)

	appController := controller.NewController(inMemoryStore, CLIView, HTTPService, settings)
	internal.LOGGER.Info("Starting the app.")
	appController.Start()
}
