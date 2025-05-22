package store

import (
	"GoHealthChecker/internal/model"
)

type Store interface {
	GetResult(url string) []model.HealthCheckResult
	GetResultsAll() map[string][]model.HealthCheckResult
	GetLatestResults() map[string]model.HealthCheckResult

	AddURL(url string) error
	SaveResult(url string, result model.HealthCheckResult)
	GetURLs() []string
}
