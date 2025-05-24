package integration

import (
	"GoHealthChecker/internal/controller"
	"GoHealthChecker/internal/service"
	"GoHealthChecker/internal/store"
	"GoHealthChecker/internal/view"
	"GoHealthChecker/tests"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestOneWorking(t *testing.T) {
	t.Parallel()
	// Enable HTTP mocking
	// Register mock responses
	httpmockTransport := httpmock.NewMockTransport()
	httpmockTransport.RegisterResponder("GET", "https://example1.com",
		httpmock.NewStringResponder(200, "<html><body>Example Domain</body></html>").Delay(100*time.Millisecond),
	)

	// the app will ping servers each second
	output, _, cancel, settings := tests.CreateConfiguration(1, 1)

	// Initialize app components
	inMemoryStore := store.NewInMemoryStore()
	cliView := view.NewCLIView(settings)
	httpService := service.NewHTTPServiceWithClient(
		&http.Client{
			Transport: httpmockTransport,
		},
	)
	appController := controller.NewController(inMemoryStore, cliView, httpService, settings)
	done := make(chan struct{})

	// Run the app in a goroutine
	go func() {
		_ = appController.Start([]string{"https://example1.com"})
		close(done) // Signal that the app has finished
	}()

	// Let it run for a short time
	// the app will ping server each second
	// this results in 1 initial + 5 seconds of pings
	time.Sleep(5*time.Second + 500*time.Millisecond)

	// Stop the app
	cancel()

	<-done
	info := httpmockTransport.GetCallCountInfo()
	contentStr := output.String()
	contentStr = strings.ReplaceAll(contentStr, "\u001B[H\u001B[2J", "")
	exampleCalls := tests.ParseLinesForURL(contentStr, "https://example1.com")

	// 6 rerenders happened
	assert.Equal(t, 7, len(exampleCalls))
	for i := 0; i < 6; i++ {
		tests.VerifyRerenders(t, exampleCalls[i], "UP", "200", 100.00, 40)
	}
	tests.VerifyMetricsTable(t, exampleCalls[6], "6/0", "100.0%", 100.00, 40, 100.00, 40, 100.00, 40)
	// Verify HTTP mock was actually called
	// only 6 calls are done, the 7th record is the last rerender with the table result
	assert.Equal(t, 6, info["GET https://example1.com"])
}

func TestTwoWorking(t *testing.T) {
	t.Parallel()

	// Enable HTTP mocking
	httpmockTransport := httpmock.NewMockTransport()
	defer httpmockTransport.Reset() // Cleanup after test

	// Register mock responses
	httpmockTransport.RegisterResponder("GET", "https://example2.com",
		httpmock.NewStringResponder(200, "<html><body>Example Domain</body></html>").Delay(100*time.Millisecond),
	)
	httpmockTransport.RegisterResponder("GET", "https://example2.org",
		httpmock.NewStringResponder(200, "<html><body>Example Domain</body></html>").Delay(100*time.Millisecond),
	)

	// the app will ping servers each second
	output, _, cancel, settings := tests.CreateConfiguration(1, 1)

	// Initialize app components
	inMemoryStore := store.NewInMemoryStore()
	cliView := view.NewCLIView(settings)
	httpService := service.NewHTTPServiceWithClient(
		&http.Client{
			Transport: httpmockTransport,
		},
	)
	appController := controller.NewController(inMemoryStore, cliView, httpService, settings)
	done := make(chan struct{})

	// Run the app in a goroutine
	go func() {
		_ = appController.Start([]string{"https://example2.com", "https://example2.org"})
		close(done)
	}()

	// Let it run for a short time
	// the app will ping server each second
	// this results in 1 initial + 5 seconds of pings
	time.Sleep(5*time.Second + 500*time.Millisecond)

	// Stop the app
	cancel()
	<-done

	contentStr := output.String()
	contentStr = strings.ReplaceAll(contentStr, "\u001B[H\u001B[2J", "")
	exampleCalls := tests.ParseLinesForURL(contentStr, "https://example2.com")
	exampleOrgCalls := tests.ParseLinesForURL(contentStr, "https://example2.org")

	// 6 rerenders happened, verify
	for i := 0; i < len(exampleCalls)-1; i++ {
		tests.VerifyRerenders(t, exampleCalls[i], "UP", "200", 100.00, 40)
	}
	// 6 rerenders happened, verify example2.org
	for i := 0; i < len(exampleOrgCalls)-1; i++ {
		tests.VerifyRerenders(t, exampleOrgCalls[i], "UP", "200", 100.00, 40)
	}
	// verify resulting metrics
	tests.VerifyMetricsTable(t, exampleCalls[len(exampleCalls)-1], "6/0", "100.0%", 100.00, 40, 100.00, 40, 100.00, 40)
	tests.VerifyMetricsTable(t, exampleOrgCalls[len(exampleOrgCalls)-1], "6/0", "100.0%", 100.00, 40, 100.00, 40, 100.00, 40)
	// Verify HTTP mock was actually called
	info := httpmockTransport.GetCallCountInfo()
	// only 6 calls are done, the 7th record is the last rerender with the table result
	assert.Equal(t, 6, info["GET https://example2.com"])
	assert.Equal(t, 6, info["GET https://example2.org"])
}

func TestOneWorkingOne500(t *testing.T) {
	t.Parallel()

	// Enable HTTP mocking
	httpmockTransport := httpmock.NewMockTransport()
	defer httpmockTransport.Reset() // Cleanup after test

	// Register mock responses
	httpmockTransport.RegisterResponder("GET", "https://example2.com",
		httpmock.NewStringResponder(200, "<html><body>Example Domain</body></html>").Delay(100*time.Millisecond),
	)
	httpmockTransport.RegisterResponder("GET", "https://error.org",
		httpmock.NewStringResponder(500, "<html><body>Example Domain</body></html>").Delay(100*time.Millisecond),
	)

	// the app will ping servers each second
	output, _, cancel, settings := tests.CreateConfiguration(1, 1)

	// Initialize app components
	inMemoryStore := store.NewInMemoryStore()
	cliView := view.NewCLIView(settings)
	httpService := service.NewHTTPServiceWithClient(
		&http.Client{
			Transport: httpmockTransport,
		},
	)
	appController := controller.NewController(inMemoryStore, cliView, httpService, settings)
	done := make(chan struct{})

	// Run the app in a goroutine
	go func() {
		_ = appController.Start([]string{"https://example2.com", "https://error.org"})
		close(done)
	}()

	// Let it run for a short time
	// the app will ping server each second
	// this results in 1 initial + 5 seconds of pings
	time.Sleep(5*time.Second + 500*time.Millisecond)

	// Stop the app
	cancel()
	<-done

	contentStr := output.String()
	contentStr = strings.ReplaceAll(contentStr, "\u001B[H\u001B[2J", "")
	exampleCalls := tests.ParseLinesForURL(contentStr, "https://example2.com")
	errorCalls := tests.ParseLinesForURL(contentStr, "https://error.org")

	// 6 rerenders happened, verify
	for i := 0; i < len(exampleCalls)-1; i++ {
		tests.VerifyRerenders(t, exampleCalls[i], "UP", "200", 100.00, 40)
	}
	// 6 rerenders happened, verify example2.org
	for i := 0; i < len(errorCalls)-1; i++ {
		tests.VerifyRerenders(t, errorCalls[i], "DOWN", "500", 100.00, 40)
	}
	// verify resulting metrics
	tests.VerifyMetricsTable(t, exampleCalls[len(exampleCalls)-1], "6/0", "100.0%", 100.00, 40, 100.00, 40, 100.00, 40)
	tests.VerifyMetricsTable(t, errorCalls[len(errorCalls)-1], "0/6", "0.0%", 100.00, 40, 100.00, 40, 100.00, 40)
	// Verify HTTP mock was actually called
	info := httpmockTransport.GetCallCountInfo()
	// only 6 calls are done, the 7th record is the last rerender with the table result
	assert.Equal(t, 6, info["GET https://example2.com"])
	assert.Equal(t, 6, info["GET https://error.org"])
}

func TestInvalidArgs(t *testing.T) {
	t.Parallel()

	_, _, _, settings := tests.CreateConfiguration(1, 1)

	// Initialize app components
	inMemoryStore := store.NewInMemoryStore()
	cliView := view.NewCLIView(settings)
	httpService := service.NewHTTPService(settings)
	appController := controller.NewController(inMemoryStore, cliView, httpService, settings)

	// Test various invalid URLs
	invalidURLs := []string{
		"htt://invalid-url",
		"://missing-scheme",
		"http:/missing-slash",
		"",                      // Empty string
		"http://invalid@host$",  // Invalid characters
		"ftp://unsupported.com", // Unsupported scheme
		"https://:443/path",
		"https://this-domain-should-not-exist-anywhere.net",
	}

	for _, url := range invalidURLs {
		t.Logf("Testing invalid URL: %s", url)
		err := appController.Start([]string{url})
		assert.Error(t, err, "Expected error for invalid URL: %s", url)
	}
}

func TestNoArgs(t *testing.T) {
	t.Parallel()

	_, _, _, settings := tests.CreateConfiguration(1, 1)

	// Initialize app components
	inMemoryStore := store.NewInMemoryStore()
	cliView := view.NewCLIView(settings)
	httpService := service.NewHTTPService(settings)
	appController := controller.NewController(inMemoryStore, cliView, httpService, settings)

	err := appController.Start([]string{})
	assert.Error(t, err, "Expected error for no URLs provided")
}

func TestOneSlowWorking(t *testing.T) {
	t.Parallel()
	// Enable HTTP mocking
	// Register mock responses
	httpmockTransport := httpmock.NewMockTransport()
	httpmockTransport.RegisterResponder("GET", "https://example1.com",
		httpmock.NewStringResponder(200, "<html><body>Example Domain</body></html>").Delay(2000*time.Millisecond),
	)

	// the app will ping servers each second
	output, _, cancel, settings := tests.CreateConfiguration(1, 1)
	settings.WithMaxQueueSize(1000) // Increase queue size to handle slow responses

	// Initialize app components
	inMemoryStore := store.NewInMemoryStore()
	cliView := view.NewCLIView(settings)
	httpService := service.NewHTTPServiceWithClient(
		&http.Client{
			Transport: httpmockTransport,
		},
	)
	appController := controller.NewController(inMemoryStore, cliView, httpService, settings)
	done := make(chan struct{})

	// Run the app in a goroutine
	go func() {
		_ = appController.Start([]string{"https://example1.com"})
		close(done) // Signal that the app has finished
	}()

	// 1 element is added every second -> 10 elements in total
	// elements are processed every 2 seconds
	// queue will have 5 elements after 10 seconds
	time.Sleep(10*time.Second + 500*time.Millisecond)

	// Stop the app
	cancel()

	<-done
	info := httpmockTransport.GetCallCountInfo()
	contentStr := output.String()
	contentStr = strings.ReplaceAll(contentStr, "\u001B[H\u001B[2J", "")
	exampleCalls := tests.ParseLinesForURL(contentStr, "https://example1.com")

	// 11 rerenders happened + 1 metrics table
	assert.Equal(t, 12, len(exampleCalls))
	for i := 0; i < 11; i++ {
		tests.VerifyRerenders(t, exampleCalls[i], "UP", "200", 2000.00, 40)
	}
	tests.VerifyMetricsTable(t, exampleCalls[11], "11/0", "100.0%", 2000.00, 40, 2000.00, 40, 2000.00, 40)
	assert.Equal(t, 11, info["GET https://example1.com"])
}

func TestHTTPClientError(t *testing.T) {
	t.Parallel()
	// Enable HTTP mocking
	// Register mock responses
	httpmockTransport := httpmock.NewMockTransport()
	httpmockTransport.RegisterResponder("GET", "https://example1.com",
		httpmock.NewErrorResponder(errors.New("connection refused")),
	)

	// the app will ping servers each second
	output, _, cancel, settings := tests.CreateConfiguration(1, 1)

	// Initialize app components
	inMemoryStore := store.NewInMemoryStore()
	cliView := view.NewCLIView(settings)
	httpService := service.NewHTTPServiceWithClient(
		&http.Client{
			Transport: httpmockTransport,
		},
	)
	appController := controller.NewController(inMemoryStore, cliView, httpService, settings)
	done := make(chan struct{})

	// Run the app in a goroutine
	go func() {
		_ = appController.Start([]string{"https://example1.com"})
		close(done) // Signal that the app has finished
	}()

	// Let it run for a short time
	// the app will ping server each second
	// this results in 1 initial + 5 seconds of pings
	time.Sleep(5*time.Second + 500*time.Millisecond)

	// Stop the app
	cancel()

	<-done
	info := httpmockTransport.GetCallCountInfo()
	contentStr := output.String()
	contentStr = strings.ReplaceAll(contentStr, "\u001B[H\u001B[2J", "")
	exampleCalls := tests.ParseLinesForURL(contentStr, "https://example1.com")

	// 6 rerenders happened
	assert.Equal(t, 7, len(exampleCalls))
	for i := 0; i < 6; i++ {
		assert.Equal(t, "DOWN", exampleCalls[i][1])
		assert.Equal(t, "ERROR", exampleCalls[i][2])
		// latency is not tested, just test its here
		assert.NotEmpty(t, exampleCalls[i][3])
		assert.Equal(t, "ERROR", exampleCalls[i][4]) // Size should be 0 for errors
		assert.NotEmpty(t, exampleCalls[i][5])
	}
	assert.Equal(t, "0/6", exampleCalls[6][1])
	assert.Equal(t, "0.0%", exampleCalls[6][2])
	assert.Equal(t, "ERROR", exampleCalls[6][3])
	assert.Equal(t, "0 B", exampleCalls[6][4])
	assert.Equal(t, "ERROR", exampleCalls[6][5])
	assert.Equal(t, "0 B", exampleCalls[6][6])
	assert.Equal(t, "ERROR", exampleCalls[6][7])
	assert.Equal(t, "0 B", exampleCalls[6][8])

	// Verify HTTP mock was actually called
	// only 6 calls are done, the 7th record is the last rerender with the table result
	assert.Equal(t, 6, info["GET https://example1.com"])
}

func TestTimeout(t *testing.T) {
	t.Parallel()
	// Enable HTTP mocking
	// Register mock responses
	httpmockTransport := httpmock.NewMockTransport()
	httpmockTransport.RegisterResponder("GET", "https://example1.com",
		httpmock.NewStringResponder(200, "<html><body>Example Domain</body></html>").
			Delay(10000*time.Millisecond),
	)

	// the app will ping servers each second
	output, _, cancel, settings := tests.CreateConfiguration(1, 1)
	settings.WithTimeout(100 * time.Millisecond)

	// Initialize app components
	inMemoryStore := store.NewInMemoryStore()
	cliView := view.NewCLIView(settings)
	httpService := service.NewHTTPServiceWithClient(
		&http.Client{
			Transport: httpmockTransport,
			Timeout:   settings.Timeout,
		})

	appController := controller.NewController(inMemoryStore, cliView, httpService, settings)
	done := make(chan struct{})

	// Run the app in a goroutine
	go func() {
		_ = appController.Start([]string{"https://example1.com"})
		close(done) // Signal that the app has finished
	}()

	// Let it run for a short time
	// the app will ping server each second
	// this results in 1 initial + 5 seconds of pings
	time.Sleep(5*time.Second + 500*time.Millisecond)

	// Stop the app
	cancel()

	<-done
	info := httpmock.GetCallCountInfo()
	contentStr := output.String()
	contentStr = strings.ReplaceAll(contentStr, "\u001B[H\u001B[2J", "")
	exampleCalls := tests.ParseLinesForURL(contentStr, "https://example1.com")

	// 6 rerenders happened
	assert.Equal(t, 7, len(exampleCalls))
	for i := 0; i < 6; i++ {
		assert.Equal(t, "DOWN", exampleCalls[i][1])
		assert.Equal(t, "ERROR", exampleCalls[i][2])
		// latency is not tested, just test its here
		assert.NotEmpty(t, exampleCalls[i][3])
		assert.Equal(t, "ERROR", exampleCalls[i][4]) // Size should be 0 for errors
		assert.NotEmpty(t, exampleCalls[i][5])
	}
	assert.Equal(t, "0/6", exampleCalls[6][1])
	assert.Equal(t, "0.0%", exampleCalls[6][2])
	assert.GreaterOrEqual(t, tests.ParseFloatFromString(exampleCalls[6][3]), 100.00*0.8)
	assert.LessOrEqual(t, tests.ParseFloatFromString(exampleCalls[6][3]), 100.00*1.2)
	assert.Equal(t, "0 B", exampleCalls[6][4])
	assert.GreaterOrEqual(t, tests.ParseFloatFromString(exampleCalls[6][5]), 100.00*0.8)
	assert.LessOrEqual(t, tests.ParseFloatFromString(exampleCalls[6][5]), 100.00*1.2)
	assert.Equal(t, "0 B", exampleCalls[6][6])
	assert.GreaterOrEqual(t, tests.ParseFloatFromString(exampleCalls[6][5]), 100.00*0.8)
	assert.LessOrEqual(t, tests.ParseFloatFromString(exampleCalls[6][5]), 100.00*1.2)
	assert.Equal(t, "0 B", exampleCalls[6][8])

	// 0 requests were successful, all failed due to timeout
	assert.Equal(t, 0, info["GET https://example1.com"])
}

func TestHTTPQueueLimit(t *testing.T) {
	/*
				HTTP Delay is - 4s
		        Timeout is - 1s
			    App runs - 10s
			    Queue size is - 1

				This means that the app will try to ping the server every second, but since the server responds after 4 seconds,
				the app will have to wait for the response before it can process the next request,
				but meanwhile that the queue will get 3 more requests in the queue.

			    Total requests -> 1 nitial, 2 more processed within the 10 seconds, one in the queue -> 4
	*/
	t.Parallel()
	httpmockTransport := httpmock.NewMockTransport()
	httpmockTransport.RegisterResponder("GET", "https://example1.com",
		httpmock.NewStringResponder(200, "<html><body>Example Domain</body></html>").
			Delay(4*time.Second),
	)

	// the app will ping servers each second
	output, _, cancel, settings := tests.CreateConfiguration(10, 1)
	settings.WithMaxQueueSize(1) // Set a small queue size to trigger queue limit

	// Initialize app components
	inMemoryStore := store.NewInMemoryStore()
	cliView := view.NewCLIView(settings)
	httpService := service.NewHTTPServiceWithClient(
		&http.Client{
			Transport: httpmockTransport,
		},
	)
	appController := controller.NewController(inMemoryStore, cliView, httpService, settings)
	done := make(chan struct{})

	// Run the app in a goroutine
	go func() {
		_ = appController.Start([]string{"https://example1.com"})
		close(done) // Signal that the app has finished
	}()

	// Let it run for a short time
	time.Sleep(10 * time.Second)

	// Stop the app
	cancel()

	<-done
	info := httpmockTransport.GetCallCountInfo()
	contentStr := output.String()
	contentStr = strings.ReplaceAll(contentStr, "\u001B[H\u001B[2J", "")
	exampleCalls := tests.ParseLinesForURL(contentStr, "https://example1.com")

	assert.GreaterOrEqual(t, len(exampleCalls), 3) // Initial + at least 2 rerenders due to queue limit

	for i := 0; i < len(exampleCalls)-1; i++ {
		tests.VerifyRerenders(t, exampleCalls[i], "UP", "200", 4000, 40)
	}
	tests.VerifyMetricsTable(t, exampleCalls[len(exampleCalls)-1], "4/0", "100.0%",
		4000, 40,
		4000, 40,
		4000, 40,
	)

	assert.Equal(t, 4, info["GET https://example1.com"])
}
