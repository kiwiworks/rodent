package cas

import (
	"github.com/pkg/errors"

	"rodent/config"
	"rodent/module"
	"rodent/system/manifest"
)

func configProvider(manifest *manifest.Manifest) (*CasdoorClientConfig, error) {
	type EnvironmentConfig struct {
		ManifestPath string `split_words:"true" required:"true"`
	}

	cfg, err := config.FromEnv[EnvironmentConfig](manifest.Application, "cas")
	if err != nil {
		return nil, errors.Wrapf(err, "unable to load config from environment")
	}
	casManifest, err := ManifestFromFile(cfg.ManifestPath)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to load manifest from path %s", cfg.ManifestPath)
	}
	return casManifest.asCasdoorConfig(cfg.ManifestPath)
}

func Module() module.Module {
	return module.New(
		"core.cas",
		module.Private(
			configProvider,
		),
		module.Public(
			NewCasdoorClient,
			NewMiddleware,
		),
	)
}
