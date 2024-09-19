package object

import (
	"github.com/pkg/errors"

	"github.com/kiwiworks/rodent/config"
	"github.com/kiwiworks/rodent/module"
	"github.com/kiwiworks/rodent/system/logger"
	"github.com/kiwiworks/rodent/system/manifest"
)

func configProvider(manifest *manifest.Manifest) (*StoreConfig, error) {
	type Environment struct {
		Endpoint  string `required:"true" split_words:"true"`
		AccessKey string `required:"true" split_words:"true"`
		SecretKey string `required:"true" split_words:"true"`
		Secure    bool   `default:"true" split_words:"true"`
	}
	env, err := config.FromEnv[Environment](manifest.Application, "object")
	if err != nil {
		return nil, errors.Wrap(err, "unable to load object store config from env")
	}

	if !env.Secure {
		logger.New().Warn("object store secure access is disabled in config")
	}

	return &StoreConfig{
		Endpoint:  env.Endpoint,
		AccessKey: env.AccessKey,
		SecretKey: env.SecretKey,
		UseSSL:    env.Secure,
	}, nil
}

func Module() module.Module {
	return module.New(
		"database.object",
		module.Private(configProvider),
		module.Public(NewStore),
	)
}
