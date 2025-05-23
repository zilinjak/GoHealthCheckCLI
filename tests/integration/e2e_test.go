package integration

import (
	"GoHealthChecker/internal/controller"
	"GoHealthChecker/internal/service"
	"GoHealthChecker/internal/store"
	"GoHealthChecker/internal/view"
	"GoHealthChecker/tests"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestOneWorking(t *testing.T) {
	// Enable HTTP mocking
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Register mock responses
	httpmock.RegisterResponder("GET", "https://example.com",
		httpmock.NewStringResponder(200, "<html><body>Example Domain</body></html>").Delay(100*time.Millisecond),
	)

	os.Args = []string{
		"",
		"https://example.com",
	}
	// the app will ping servers each second
	output, _, cancel, settings := tests.CreateConfiguration(1, 1)

	// Initialize app components
	inMemoryStore := store.NewInMemoryStore()
	cliView := view.NewCLIView(settings)
	httpService := service.NewHTTPService(settings)
	appController := controller.NewController(inMemoryStore, cliView, httpService, settings)

	// Run the app in a goroutine
	go func() {
		_ = appController.Start()
	}()

	// Let it run for a short time
	// the app will ping server each second
	// this results in 1 initial + 5 seconds of pings
	time.Sleep(5*time.Second + 500*time.Millisecond)

	// Stop the app
	cancel()

	// Wait for the app to finish
	time.Sleep(2 * time.Second)

	contentStr := output.String()
	contentStr = strings.ReplaceAll(contentStr, "\u001B[H\u001B[2J", "")
	exampleCalls := tests.ParseLinesForURL(contentStr, "https://example.com")

	// 6 rerenders happened
	assert.Equal(t, 7, len(exampleCalls))
	for i := 0; i < 6; i++ {
		tests.VerifyRerenders(t, exampleCalls[i], "UP", "200", 100.00, 40)
	}
	tests.VerifyMetricsTable(t, exampleCalls[6], "6/0", "100.0%", 100.00, 40, 100.00, 40, 100.00, 40)
	// Verify HTTP mock was actually called
	info := httpmock.GetCallCountInfo()
	// only 6 calls are done, the 7th record is the last rerender with the table result
	assert.Equal(t, 6, info["GET https://example.com"])
}

func TestTwoWorking(t *testing.T) {
	// Enable HTTP mocking
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Register mock responses
	httpmock.RegisterResponder("GET", "https://example.com",
		httpmock.NewStringResponder(200, "<html><body>Example Domain</body></html>").Delay(100*time.Millisecond),
	)
	httpmock.RegisterResponder("GET", "https://example.org",
		httpmock.NewStringResponder(200, "<html><body>Example Domain</body></html>").Delay(100*time.Millisecond),
	)

	os.Args = []string{
		"",
		"https://example.com",
		"https://example.org",
	}
	// the app will ping servers each second
	output, _, cancel, settings := tests.CreateConfiguration(1, 1)

	// Initialize app components
	inMemoryStore := store.NewInMemoryStore()
	cliView := view.NewCLIView(settings)
	httpService := service.NewHTTPService(settings)
	appController := controller.NewController(inMemoryStore, cliView, httpService, settings)

	// Run the app in a goroutine
	go func() {
		_ = appController.Start()
	}()

	// Let it run for a short time
	// the app will ping server each second
	// this results in 1 initial + 5 seconds of pings
	time.Sleep(5*time.Second + 500*time.Millisecond)

	// Stop the app
	cancel()

	// Wait for the app to finish
	time.Sleep(2 * time.Second)

	contentStr := output.String()
	contentStr = strings.ReplaceAll(contentStr, "\u001B[H\u001B[2J", "")
	exampleCalls := tests.ParseLinesForURL(contentStr, "https://example.com")
	exampleOrgCalls := tests.ParseLinesForURL(contentStr, "https://example.com")

	// 6 rerenders happened
	for i := 0; i < len(exampleCalls)-1; i++ {
		tests.VerifyRerenders(t, exampleCalls[i], "UP", "200", 100.00, 40)
	}
	for i := 0; i < len(exampleOrgCalls)-1; i++ {
		tests.VerifyRerenders(t, exampleOrgCalls[i], "UP", "200", 100.00, 40)
	}

	tests.VerifyMetricsTable(t, exampleCalls[len(exampleCalls)-1], "6/0", "100.0%", 100.00, 40, 100.00, 40, 100.00, 40)
	tests.VerifyMetricsTable(t, exampleOrgCalls[len(exampleOrgCalls)-1], "6/0", "100.0%", 100.00, 40, 100.00, 40, 100.00, 40)
	// Verify HTTP mock was actually called
	info := httpmock.GetCallCountInfo()
	// only 6 calls are done, the 7th record is the last rerender with the table result
	assert.Equal(t, 6, info["GET https://example.com"])
	assert.Equal(t, 6, info["GET https://example.org"])
}
