package model

import (
	"context"
	"io"
)

type AppSettings struct {
	Timeout         int
	PollingInterval int
	Context         context.Context
	OutputStream    io.Writer
}
