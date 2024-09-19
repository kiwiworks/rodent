package migration

import (
	"context"
	"sort"

	"go.uber.org/fx"
)

type Migrator struct {
	migrations []*Migration
}

func (m *Migrator) sortMigrationsByTime() {
	sort.SliceStable(m.migrations, func(i, j int) bool {
		return m.migrations[i].time.Before(m.migrations[j].time)
	})
}

type Config struct {
	fx.In
	Migrations []*Migration `group:"migration.migration"`
}

func NewMigrator(cfg Config) *Migrator {
	m := &Migrator{
		migrations: cfg.Migrations,
	}
	m.sortMigrationsByTime()
	return m
}

func (m *Migrator) MigrateUp(ctx context.Context) error {
	for _, migration := range m.migrations {
		if err := migration.up(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (m *Migrator) MigrateDown(ctx context.Context) error {
	for _, migration := range m.migrations {
		if err := migration.down(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (m *Migrator) OnStart(ctx context.Context) error {
	return m.MigrateUp(ctx)
}
