package store

import (
	"GoHealthChecker/internal/model"
	"errors"
	urllib "net/url"
	"regexp"
)

type Store interface {
	AddURL(url string) error
	SaveResult(url string, result model.HealthCheckResult)

	GetURLs() []string
	GetLatestResults() map[string]model.HealthCheckResult
	GetMetrics() map[string]model.Metrics
}

func ValidateURL(url string) error {
	// This method is not implemented in the interface, but can be used by implementations
	// check if the URL is valid
	item, err := urllib.ParseRequestURI(url)
	if err != nil {
		return errors.New("invalid URL: " + err.Error())
	}

	if item.Scheme != "http" && item.Scheme != "https" {
		return errors.New("unsupported URL scheme: " + item.Scheme)
	}
	// Check host - must have a non-empty hostname
	if item.Host == "" {
		return errors.New("URL must have a host")
	}

	// Validate hostname format
	hostname := item.Hostname()
	if hostname == "" {
		return errors.New("invalid hostname")
	}

	// Check for invalid characters in hostname using regex
	// Valid hostname: letters, digits, hyphens, dots
	hostnameRegex := regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9\-\.]+[a-zA-Z0-9]$`)
	if !hostnameRegex.MatchString(hostname) {
		return errors.New("hostname contains invalid characters")
	}

	//_, err = net.LookupHost(hostname)
	//if err != nil {
	//	return fmt.Errorf("DNS resolution failed: %w", err)
	//}
	return nil
}
