package controller

import (
	"GoHealthChecker/internal"
	"GoHealthChecker/internal/model"
	"GoHealthChecker/internal/service"
	"GoHealthChecker/internal/store"
	"GoHealthChecker/internal/view"
	"context"
	"fmt"
	"os"
	"sync"
	"time"
)

type Controller struct {
	HTTPService    service.Service
	Store          store.Store
	View           view.View
	workersWg      sync.WaitGroup
	workerChannels map[string]chan struct{}
	channelsMutex  sync.RWMutex
	settings       model.AppSettings
}

func NewController(
	store store.Store,
	view view.View,
	service service.Service,
	settings model.AppSettings,
) *Controller {
	return &Controller{
		HTTPService:    service,
		Store:          store,
		View:           view,
		workersWg:      sync.WaitGroup{},
		workerChannels: make(map[string]chan struct{}),
		channelsMutex:  sync.RWMutex{},
		settings:       settings,
	}
}

func (controller *Controller) Start(urls []string) error {
	// Parse args and load them to Store
	err := controller.ParseArgs(urls)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return err
	}

	ticker := time.NewTicker(time.Duration(controller.settings.PollingInterval) * time.Second)
	defer ticker.Stop()

	internal.LOGGER.Info("Starting loop (press Ctrl+C to stop)...")
	// Add elements to Queues
	internal.LOGGER.Info("Spawning workers for each URL")
	_, cancel := controller.startWorkers()
	controller.addToQueue()
	for {
		select {
		case <-controller.settings.Context.Done():
			internal.LOGGER.Info("Gracefully exiting...")
			cancel() // Cancel worker context
			controller.Stop()
			return nil
		case <-ticker.C:
			// every 5 seconds we add to the channels
			controller.addToQueue()
		}
	}
}

func (controller *Controller) Stop() {
	// Wait for all workers to finish
	controller.workersWg.Wait()
	controller.View.RenderMetrics(controller.Store.GetMetrics())
}

func (controller *Controller) ParseArgs(urls []string) error {
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
	controller.channelsMutex.Lock()
	defer controller.channelsMutex.Unlock()

	// Add elements to Queues
	for _, url := range controller.Store.GetURLs() {
		if _, exists := controller.workerChannels[url]; !exists {
			controller.workerChannels[url] = make(chan struct{}, 5)
		}
		if len(controller.workerChannels[url]) == cap(controller.workerChannels[url]) {
			internal.LOGGER.Warn(fmt.Sprintf("Queue for %s is FULL\n", url))
			continue
		}
		controller.workerChannels[url] <- struct{}{}
	}
	//fmt.Println("Currently these are the queues and their size:")
	for url, ch := range controller.workerChannels {
		internal.LOGGER.Info(fmt.Sprintf("URL: %s, Queue size: %d\n", url, len(ch)))
	}
}

func (controller *Controller) worker(url string, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			internal.LOGGER.Info(fmt.Sprintf("Worker for %s stopped\n", url))
			controller.workersWg.Done()
			return
		case <-func() chan struct{} {
			controller.channelsMutex.RLock()
			ch := controller.workerChannels[url]
			controller.channelsMutex.RUnlock()
			return ch
		}():
			resp, err := controller.HTTPService.CheckUrl(url)
			if err != nil {
				internal.LOGGER.Error(fmt.Sprintf("Error when requesting %s: %s", url, err))
			}
			controller.Store.SaveResult(url, resp)
			controller.View.Render(controller.Store.GetLatestResults())
		}
	}
}

func (controller *Controller) startWorkers() (context.Context, context.CancelFunc) {
	// Add elements to Queues
	internal.LOGGER.Info("Spawning workers for each URL")
	workerCtx, cancel := context.WithCancel(controller.settings.Context)
	for _, url := range controller.Store.GetURLs() {
		controller.workersWg.Add(1)
		go controller.worker(url, workerCtx)
	}
	return workerCtx, cancel
}
