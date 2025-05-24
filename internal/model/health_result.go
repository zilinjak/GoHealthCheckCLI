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
	Error      error         `json:"-"`           // Error if any occurred during the check
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

func NewHealthCheckResultWithError(err error, latency time.Duration) HealthCheckResult {
	return HealthCheckResult{
		IsOk:      false,
		Latency:   latency,
		Error:     err,
		Timestamp: time.Now().UTC(),
	}
}
