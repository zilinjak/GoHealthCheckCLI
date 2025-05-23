package service

import (
	"GoHealthChecker/internal"
	"GoHealthChecker/internal/model"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Service interface {
	CheckUrl(url string) (model.HealthCheckResult, error)
}

type HTTPService struct {
	client *http.Client
}

func NewHTTPService(settings model.AppSettings) *HTTPService {
	return &HTTPService{
		client: &http.Client{
			// disable follow redirects
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Timeout: time.Duration(settings.Timeout) * time.Second, // Default timeout
		},
	}
}

func (H HTTPService) CheckUrl(url string) (model.HealthCheckResult, error) {
	start := time.Now()
	resp, err := H.client.Get(url)
	duration := time.Since(start)
	internal.LOGGER.Info(fmt.Sprintf("%s: %s\n", url, duration.String()))

	if err != nil {
		return model.HealthCheckResult{}, err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.HealthCheckResult{}, err
	}
	var sizeOfResponse = uint64(len(data))
	return model.NewHealthCheckResult(resp.StatusCode, duration, sizeOfResponse), nil
}
