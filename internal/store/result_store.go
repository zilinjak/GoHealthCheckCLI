package store

import (
	"GoHealthChecker/internal/model"
	"sync"
)

type ResultRepository interface {
	Save(result model.HealthCheckResult)
	GetResult(url string) []model.HealthCheckResult
	GetResultsAll() map[string][]model.HealthCheckResult
}

type InMemoryStore struct {
	mu    sync.RWMutex
	items map[string][]model.HealthCheckResult // Map of URL to result history
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		items: make(map[string][]model.HealthCheckResult),
	}
}

func (s *InMemoryStore) Save(result model.HealthCheckResult) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items[result.URL] = append(s.items[result.URL], result)
}

func (s *InMemoryStore) GetResult(url string) []model.HealthCheckResult {
	s.mu.RLock()
	defer s.mu.RUnlock()

	history := make([]model.HealthCheckResult, len(s.items[url]))
	copy(history, s.items[url])
	return history
}

func (s *InMemoryStore) GetResultsAll() map[string][]model.HealthCheckResult {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// TODO: This can maybe be optimzed
	resultsCopy := make(map[string][]model.HealthCheckResult)
	for k, v := range s.items {
		resultsCopy[k] = make([]model.HealthCheckResult, len(v))
		copy(resultsCopy[k], v)
	}
	return resultsCopy
}
