package nosql

import (
	"net/url"

	"github.com/kiwiworks/rodent/config"
	"github.com/kiwiworks/rodent/module"
	"github.com/kiwiworks/rodent/system/errors"
	"github.com/kiwiworks/rodent/system/manifest"
)

func configProvider(manifest *manifest.Manifest) (*SurrealDBConfig, error) {
	type Environment struct {
		Dsn url.URL `required:"true"`
	}
	env, err := config.FromEnv[Environment](manifest.Application, "surrealdb")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load environment")
	}
	return ConfigFromUrl(env.Dsn)
}

func Module() module.Module {
	return module.New(
		"database.nosql",
		module.Private(configProvider),
		module.Public(NewClient),
		module.Service[Client](),
	)
}
