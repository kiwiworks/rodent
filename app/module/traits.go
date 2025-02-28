package module

import (
	"context"
	"time"
)

type HealthCheckManifest struct {
	Name        string
	Description string
	StartedAt   time.Time
	CrashedAt   *time.Time
	LastError   error
}

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
	HealthCheck interface {
		OnStartStop
		Inspect() HealthCheckManifest
	}
)
