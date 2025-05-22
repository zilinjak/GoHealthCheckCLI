package controller

import (
	"GoHealthChecker/internal/service"
	"GoHealthChecker/internal/store"
	"GoHealthChecker/internal/view"
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"
)

type Controller struct {
	HTTPService    service.Service
	Store          store.Store
	View           view.View
	workersWg      sync.WaitGroup
	workerChannels map[string]chan struct{}
}

func NewController(
	store store.Store,
	view view.View,
	service service.Service,
) *Controller {
	return &Controller{
		HTTPService:    service,
		Store:          store,
		View:           view,
		workersWg:      sync.WaitGroup{},
		workerChannels: make(map[string]chan struct{}),
	}
}

func (controller *Controller) handleError(err error) {
	// print to stderr
	fmt.Fprintln(os.Stderr, "Error:", err)
	os.Exit(1)
}

func (controller *Controller) Start() {
	// Parse args and load them to Store
	err := controller.ParseArgs()
	if err != nil {
		controller.handleError(err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	// Main loop
	fmt.Println("Starting loop (press Ctrl+C to stop)...")
	// Add elements to Queues
	fmt.Println("Spawning workers for each URL")
	ctx, cancel := context.WithCancel(context.Background())
	for _, url := range controller.Store.GetURLs() {
		controller.workersWg.Add(1)
		go controller.worker(url, ctx)
	}
	// initial queue filling
	controller.addToQueue()
	for {
		select {
		case <-stop:
			fmt.Println("\nGracefully exiting...")
			cancel()
			controller.Stop()
		case <-ticker.C:
			// every 5 seconds we add to the channels
			controller.addToQueue()
		}
	}
}

func (controller *Controller) Stop() {
	// Wait for all workers to finish
	controller.workersWg.Wait()
	controller.View.RenderTable(controller.Store.GetResultsAll())
	os.Exit(0)
}

func (controller *Controller) ParseArgs() error {
	urls := os.Args[1:]

	if len(urls) == 0 {
		return fmt.Errorf("no URLs provided")
	}

	for _, url := range urls {
		if err := controller.Store.AddURL(url); err != nil {
			return fmt.Errorf("failed to add URL %s: %w", url, err)
		}
	}
	return nil
}

func (controller *Controller) addToQueue() {
	// Add elements to Queues
	for _, url := range controller.Store.GetURLs() {
		if _, exists := controller.workerChannels[url]; !exists {
			controller.workerChannels[url] = make(chan struct{}, 5)
			// We cap the queue at max 5 requests
		}
		if len(controller.workerChannels[url]) == cap(controller.workerChannels[url]) {
			fmt.Printf("%s FULL\n", url)
			continue
		}
		controller.workerChannels[url] <- struct{}{}
	}
	//fmt.Println("Currently these are the queues and their size:")
	for url, ch := range controller.workerChannels {
		fmt.Printf("URL: %s, Queue size: %d\n", url, len(ch))
	}
}

func (controller *Controller) worker(url string, ctx context.Context) {
	// consume from the channel
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker for %s stopped\n", url)
			controller.workersWg.Done()
			return
		case <-controller.workerChannels[url]:
			resp, err := controller.HTTPService.CheckUrl(url)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error when requesting %s: %s", url, err)
			}
			controller.Store.SaveResult(url, *resp)
		}

	}
}
