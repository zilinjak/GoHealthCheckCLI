package model

import (
	"time"
)

type HealthCheckResult struct {
	StatusCode int           `json:"status_code"` // HTTP status code (0 if network error)
	Latency    time.Duration `json:"latency"`     // Request duration
	Timestamp  time.Time     `json:"timestamp"`   // When check occurred
	isOk       bool          `json:"isOk"`        // Is the URL healthy
}

func NewHealthCheckResult(
	statusCode int,
	latency time.Duration,
) *HealthCheckResult {
	isOk := statusCode >= 200 && statusCode < 400

	return &HealthCheckResult{
		isOk:       isOk,
		StatusCode: statusCode,
		Latency:    latency,
		Timestamp:  time.Now().UTC(),
	}
}
