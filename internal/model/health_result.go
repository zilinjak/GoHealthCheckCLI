package model

import (
	"time"
)

// HealthCheckResult represents the outcome of a single health check
type HealthCheckResult struct {
	URL        string        `json:"url"`         // Checked URL
	Status     string        `json:"status"`      // "UP", "DOWN", or "TIMEOUT"
	StatusCode int           `json:"status_code"` // HTTP status code (0 if network error)
	Latency    time.Duration `json:"latency"`     // Request duration
	Timestamp  time.Time     `json:"timestamp"`   // When check occurred
}

// NewHealthCheckResult constructor for safe initialization
func NewHealthCheckResult(
	url string,
	status string,
	statusCode int,
	latency time.Duration,
) HealthCheckResult {
	return HealthCheckResult{
		URL:        url,
		Status:     status,
		StatusCode: statusCode,
		Latency:    latency,
		Timestamp:  time.Now().UTC(),
	}
}
