package view

import "GoHealthChecker/internal/model"

type CLIView struct{}

func NewCLIView() *CLIView {
	return &CLIView{}
}

func (v *CLIView) Render(results map[string]model.HealthCheckResult) {
	println("View.Render")
}

func (v *CLIView) RenderTable(results map[string][]model.HealthCheckResult) {
	println("View.RenderTable")
}
