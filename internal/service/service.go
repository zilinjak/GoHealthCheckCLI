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
			Timeout: settings.Timeout, // Default timeout
		},
	}
}

func NewHTTPServiceWithClient(client *http.Client) *HTTPService {
	return &HTTPService{
		client: client,
	}
}

func (H HTTPService) CheckUrl(url string) (model.HealthCheckResult, error) {
	internal.LOGGER.Info(fmt.Sprintf("Checking %s\n", url))
	start := time.Now()
	resp, err := H.client.Get(url)
	duration := time.Since(start)

	if err != nil {
		return model.NewHealthCheckResultWithError(err, duration), err
	}
	internal.LOGGER.Info(fmt.Sprintf("%s -> %d: %s\n", url, resp.StatusCode, duration.String()))

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.NewHealthCheckResultWithError(err, duration), err
	}
	var sizeOfResponse = uint64(len(data))
	return model.NewHealthCheckResult(resp.StatusCode, duration, sizeOfResponse), nil
}
