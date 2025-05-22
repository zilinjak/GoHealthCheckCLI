package store

import (
	"GoHealthChecker/internal/model"
	urllib "net/url"
	"sync"
)

type InMemoryStore struct {
	mu    sync.RWMutex
	items map[string][]model.HealthCheckResult // Map of URL to result history
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		items: make(map[string][]model.HealthCheckResult),
	}
}

func (s *InMemoryStore) SaveResult(url string, result model.HealthCheckResult) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items[url] = append(s.items[url], result)
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

func (s *InMemoryStore) AddURL(url string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// check if the URL is valid
	if _, err := urllib.ParseRequestURI(url); err != nil {
		return err // Invalid URL
	}

	// URL already exists
	if _, exists := s.items[url]; exists {
		return nil
	}

	s.items[url] = []model.HealthCheckResult{}
	return nil
}

func (s *InMemoryStore) GetURLs() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	urls := make([]string, 0, len(s.items))
	for url := range s.items {
		urls = append(urls, url)
	}
	return urls
}

func (s *InMemoryStore) GetLatestResults() map[string]model.HealthCheckResult {
	s.mu.RLock()
	defer s.mu.RUnlock()

	latestResults := make(map[string]model.HealthCheckResult)
	for url, results := range s.items {
		if len(results) > 0 {
			latestResults[url] = results[len(results)-1]
		}
	}
	return latestResults
}
