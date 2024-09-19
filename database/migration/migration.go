package migration

import (
	"context"
	"time"
)

type (
	Up        func(context.Context) error
	Down      func(context.Context) error
	Migration struct {
		name string
		time time.Time
		up   func(ctx context.Context) error
		down func(ctx context.Context) error
	}
)

func New(name string, migrationTime time.Time, up Up, down Down) *Migration {
	return &Migration{
		name: name,
		time: migrationTime,
		up:   up,
		down: down,
	}
}
