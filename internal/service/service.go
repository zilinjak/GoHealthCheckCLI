package service

import (
	"GoHealthChecker/internal/model"
	"fmt"
	"net/http"
	"time"
)

type Service interface {
	CheckUrl(url string) (*model.HealthCheckResult, error)
}

type HTTPService struct {
	client *http.Client
}

func NewHTTPService(timeout int) *HTTPService {
	return &HTTPService{
		client: &http.Client{
			// disable follow redirects
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Timeout: time.Duration(timeout) * time.Second, // Default timeout
		},
	}
}

func (H HTTPService) CheckUrl(url string) (*model.HealthCheckResult, error) {
	start := time.Now()
	resp, err := H.client.Get(url)
	duration := time.Since(start)
	fmt.Printf("%s: %s\n", url, duration.String())
	if err != nil {
		return nil, err
	}
	return model.NewHealthCheckResult(resp.StatusCode, duration), nil
}
