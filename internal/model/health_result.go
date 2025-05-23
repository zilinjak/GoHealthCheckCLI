package model

import (
	"time"
)

type HealthCheckResult struct {
	StatusCode int           `json:"status_code"` // HTTP status code (0 if network error)
	Latency    time.Duration `json:"latency"`     // Request duration
	Timestamp  time.Time     `json:"timestamp"`   // When check occurred
	IsOk       bool          `json:"isOk"`        // Is the URL healthy
	Size       uint64        `json:"size"`        // Size of the response
}

func NewHealthCheckResult(
	statusCode int,
	latency time.Duration,
	sizeOfResponse uint64,
) HealthCheckResult {
	isOk := statusCode >= 200 && statusCode < 400

	return HealthCheckResult{
		IsOk:       isOk,
		StatusCode: statusCode,
		Latency:    latency,
		Timestamp:  time.Now().UTC(),
		Size:       sizeOfResponse,
	}
}
