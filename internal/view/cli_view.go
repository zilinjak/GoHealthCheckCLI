package view

import (
	"GoHealthChecker/internal/model"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/jedib0t/go-pretty/table"
)

// CLIView TODO: Change current implementation to use better terminal GUI library
type CLIView struct {
	output io.Writer
}

func NewCLIView(appSettings model.AppSettings) *CLIView {
	instance := &CLIView{
		output: appSettings.OutputStream,
	}
	return instance
}

func (v *CLIView) Render(results map[string]model.HealthCheckResult) {
	v.clearTerminal()

	t := table.NewWriter()
	t.SetOutputMirror(v.output)
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
		state := "UP"
		if !result.IsOk {
			state = "DOWN"
		}
		t.AppendRow(
			table.Row{
				url,
				state,
				result.StatusCode,
				result.Latency.String(),
				formatBytes(result.Size),
				result.Timestamp,
			},
		)
	}
	t.Render()
}

func (v *CLIView) RenderMetrics(results map[string]model.Metrics) {
	v.clearTerminal()

	t := table.NewWriter()
	t.SetOutputMirror(v.output)
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
				addSuffix(result.LatencyAverage, "ms"), formatBytes(result.SizeAverage),
				addSuffix(result.LatencyMin, "ms"), formatBytes(result.SizeMin),
				addSuffix(result.LatencyMax, "ms"), formatBytes(result.SizeMax),
			},
		)
	}
	t.Render()
}

func (v *CLIView) clearTerminal() {
	_, _ = fmt.Fprint(v.output, "\033[H\033[2J")
}

func addSuffix(data float64, suffix string) string {
	return fmt.Sprintf("%.2f%s", data, suffix)
}

func formatBytes(bytes uint64) string {
	if bytes == 0 {
		return "0 B"
	}

	units := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	value := float64(bytes)
	unitIndex := 0

	for value >= 1024 && unitIndex < len(units)-1 {
		value /= 1024
		unitIndex++
	}

	// Format with two decimal places, then trim trailing zeros
	strValue := fmt.Sprintf("%.2f", value)
	trimmed := strings.TrimRight(strValue, "0")
	if len(trimmed) > 0 && trimmed[len(trimmed)-1] == '.' {
		trimmed = trimmed[:len(trimmed)-1]
	}

	return fmt.Sprintf("%s %s", trimmed, units[unitIndex])
}
