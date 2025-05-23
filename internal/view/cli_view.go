package view

import (
	"GoHealthChecker/internal/model"
	"fmt"
	"os"
	"sort"

	"github.com/jedib0t/go-pretty/table"
)

// CLIView TODO: Change current implementation to use better terminal GUI library
type CLIView struct {
}

func NewCLIView() *CLIView {
	instance := &CLIView{}
	return instance
}

func (v *CLIView) Render(results map[string]model.HealthCheckResult) {
	clearTerminal()

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"URL", "Status", "StatusCode", "Latency", "Size", "Timestamp"})

	// Extract and sort the URLs
	urls := make([]string, 0, len(results))
	for url := range results {
		urls = append(urls, url)
	}
	sort.Strings(urls)

	// Iterate through sorted URLs
	for _, url := range urls {
		result := results[url]
		t.AppendRow(
			table.Row{
				url,
				result.IsOk,
				result.StatusCode,
				result.Latency.String(),
				addSuffixUint(result.Size, "B"),
				result.Timestamp,
			},
		)
	}
	t.Render()
}

func (v *CLIView) RenderMetrics(results map[string]model.Metrics) {
	clearTerminal()

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{
		"URL", "Success/Failed", "Uptime",
		"Avg. Latency", "Avg. Size",
		"Min Latency", "Min Size",
		"Max Latency", "Max Size",
	})

	// Extract and sort the URLs
	urls := make([]string, 0, len(results))
	for url := range results {
		urls = append(urls, url)
	}
	sort.Strings(urls)

	// Iterate through sorted URLs
	for _, url := range urls {
		result := results[url]
		uptime := fmt.Sprintf("%.1f%%", float64(result.SuccessRequests)/float64(result.TotalRequests)*100)
		t.AppendRow(
			table.Row{
				url,
				fmt.Sprintf("%d/%d", result.SuccessRequests, result.FailedRequests),
				uptime,
				addSuffix(result.LatencyAverage, "ms"), addSuffixUint(result.SizeAverage, "B"),
				addSuffix(result.LatencyMin, "ms"), addSuffixUint(result.SizeMin, "B"),
				addSuffix(result.LatencyMax, "ms"), addSuffixUint(result.SizeMax, "B"),
			},
		)
	}
	t.Render()
}

func clearTerminal() {
	fmt.Print("\033[H\033[2J")
}

func addSuffix(data float64, suffix string) string {
	return fmt.Sprintf("%.2f%s", data, suffix)
}

func addSuffixUint(data uint64, suffix string) string {
	return fmt.Sprintf("%d%s", data, suffix)
}
