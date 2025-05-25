package store

import (
	"GoHealthChecker/internal/model"
	"sync"
)

type InMemoryStore struct {
	mu            sync.RWMutex
	latestResults map[string]model.HealthCheckResult
	resultMetrics map[string]model.Metrics

	registeredURLs []string
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		latestResults:  make(map[string]model.HealthCheckResult),
		resultMetrics:  make(map[string]model.Metrics),
		registeredURLs: make([]string, 0),
	}
}

func (s *InMemoryStore) SaveResult(url string, result model.HealthCheckResult) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// create the latestResult && metrics incase of not initialized yet
	s.latestResults[url] = result

	if _, exists := s.resultMetrics[url]; !exists {
		// Create new metrics for first result
		s.resultMetrics[url] = model.NewMetrics(result)
	} else {
		// if already initialized then update the metrics
		metrics := s.resultMetrics[url]
		metrics.Update(result)
		s.resultMetrics[url] = metrics
	}
}

func (s *InMemoryStore) AddURL(url string) error {
	err := ValidateURL(url)
	if err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	// URL already exists in registeredURLs
	for _, registeredURL := range s.registeredURLs {
		if registeredURL == url {
			return nil
		}
	}

	// Add the URL to registeredURLs
	s.registeredURLs = append(s.registeredURLs, url)
	return nil
}

func (s *InMemoryStore) GetURLs() []string {
	return s.registeredURLs
}

func (s *InMemoryStore) GetLatestResults() map[string]model.HealthCheckResult {
	s.mu.RLock()
	defer s.mu.RUnlock()

	results := make(map[string]model.HealthCheckResult, len(s.latestResults))
	for k, v := range s.latestResults {
		results[k] = v
	}
	return results
}

func (s *InMemoryStore) GetMetrics() map[string]model.Metrics {
	s.mu.RLock()
	defer s.mu.RUnlock()

	metrics := make(map[string]model.Metrics, len(s.resultMetrics))
	for k, v := range s.resultMetrics {
		metrics[k] = v
	}
	return metrics
}
