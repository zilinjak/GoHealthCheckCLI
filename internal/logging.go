package internal

import "go.uber.org/zap"

var LOGGER *zap.Logger = NewLogger()

func NewLogger() *zap.Logger {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		"./app.log",
	}
	instance, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return instance
}
