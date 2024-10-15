package module

import (
	"context"
)

type (
	OnStart interface {
		OnStart(ctx context.Context) error
	}
	OnStop interface {
		OnStop(ctx context.Context) error
	}
	OnStartStop interface {
		OnStart
		OnStop
	}
)
