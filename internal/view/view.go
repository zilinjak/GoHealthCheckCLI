package view

import "GoHealthChecker/internal/model"

type View interface {
	Render(map[string]model.HealthCheckResult)
	RenderTable(map[string][]model.HealthCheckResult)
}
