package view

import "GoHealthChecker/internal/model"

type View interface {
	Render(map[string]model.HealthCheckResult)
	RenderMetrics(map[string]model.Metrics)
}
