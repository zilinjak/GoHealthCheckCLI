package store

import (
	"GoHealthChecker/internal/model"
)

type Store interface {
	AddURL(url string) error
	SaveResult(url string, result model.HealthCheckResult)

	GetURLs() []string
	GetLatestResults() map[string]model.HealthCheckResult
	GetMetrics() map[string]model.Metrics
}
