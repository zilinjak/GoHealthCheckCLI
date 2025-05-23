package main

import (
	"GoHealthChecker/internal"
	"GoHealthChecker/internal/controller"
	"GoHealthChecker/internal/service"
	"GoHealthChecker/internal/store"
	"GoHealthChecker/internal/view"
	"context"
	"os"
	"os/signal"
)

func main() {
	TIMEOUT := 10
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

	inMemoryStore := store.NewInMemoryStore()
	CLIView := view.NewCLIView()
	HTTPService := service.NewHTTPService(TIMEOUT)

	appController := controller.NewController(inMemoryStore, CLIView, HTTPService, ctx)
	internal.LOGGER.Info("Starting the app.")
	appController.Start()
}
