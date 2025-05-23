package internal

import "go.uber.org/zap"

// TODO: Not exactly sure if this is the best way to create logger, it introduces global state
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
