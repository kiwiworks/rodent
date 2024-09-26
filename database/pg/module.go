package pg

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/kiwiworks/rodent/config"
	"github.com/kiwiworks/rodent/errors"
	"github.com/kiwiworks/rodent/module"
	"github.com/kiwiworks/rodent/system/manifest"
)

func datasourceProvider(manifest *manifest.Manifest) (*Datasource, error) {
	type EnvironmentConfig struct {
		Dsn string `split_words:"true"`
	}
	env, err := config.FromEnv[EnvironmentConfig](manifest.Application, "postgres")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load environment config")
	}
	dsn, err := ParseDatasource(env.Dsn)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse Postgres DSN")
	}
	return dsn, nil
}

func postgresqlProvider(source *Datasource) (*Database, error) {
	dsn := source.DSN()
	db, err := sql.Open(Dialect, dsn)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to establish a connection to the postgresql database '%s'", dsn)
	}

	db.SetConnMaxLifetime(source.ConnectionSettings.MaxLifetime)
	db.SetMaxIdleConns(source.ConnectionSettings.MaxIdle)
	db.SetMaxOpenConns(source.ConnectionSettings.MaxOpen)

	return NewDatabase(db), nil
}

func Module() module.Module {
	return module.New(
		"database.pg",
		module.Private(datasourceProvider),
		module.Public(postgresqlProvider),
	)
}
