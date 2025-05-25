package model

import (
	"context"
	"io"
	"time"
)

type AppSettings struct {
	Timeout         time.Duration
	PollingInterval time.Duration
	Context         context.Context
	OutputStream    io.Writer
	MaxQueueSize    int
}

func NewAppSettings() *AppSettings {
	return &AppSettings{
		Timeout:         10 * time.Second, // default timeout
		PollingInterval: 5 * time.Second,  // default polling interval
		MaxQueueSize:    5,                // default max queue size
	}
}

func (s *AppSettings) WithTimeout(timeout time.Duration) *AppSettings {
	s.Timeout = timeout
	return s
}

func (s *AppSettings) WithPollingInterval(interval time.Duration) *AppSettings {
	s.PollingInterval = interval
	return s
}

func (s *AppSettings) WithContext(ctx context.Context) *AppSettings {
	s.Context = ctx
	return s
}

func (s *AppSettings) WithOutputStream(output io.Writer) *AppSettings {
	s.OutputStream = output
	return s
}

func (s *AppSettings) WithMaxQueueSize(size int) *AppSettings {
	s.MaxQueueSize = size
	return s
}
