package model

import "context"

type AppSettings struct {
	Timeout         int
	PoolingInterval int
	Context         context.Context
}
