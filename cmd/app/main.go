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
	// Initiaize the context and signal handler for CTRL+C handling
	ctx, _ := signalHandler()

	// Set up the application settings and components
	settings := model.AppSettings{
		Timeout:         10,
		PollingInterval: 5,
		Context:         ctx,
		OutputStream:    os.Stdout,
		MaxQueueSize:    5,
	}

	inMemoryStore := store.NewInMemoryStore()
	CLIView := view.NewCLIView(settings)
	HTTPService := service.NewHTTPService(settings)
	appController := controller.NewController(inMemoryStore, CLIView, HTTPService, settings)
	// Handle failure of the app controller - eg invalid inputs etc.
	internal.LOGGER.Info("Starting the app.")
	err := appController.Start(os.Args[1:])
	if err != nil {
		internal.LOGGER.Error("Error starting the app:" + err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}
