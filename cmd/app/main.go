package main

import (
	"GoHealthChecker/internal/controller"
	"GoHealthChecker/internal/service"
	"GoHealthChecker/internal/store"
	"GoHealthChecker/internal/view"
)

func main() {
	TIMEOUT := 10

	inMemoryStore := store.NewInMemoryStore()
	CLIView := view.NewCLIView()
	HTTPService := service.NewHTTPService(TIMEOUT)

	appController := controller.NewController(inMemoryStore, CLIView, HTTPService)

	appController.Start()
}
