package config

import (
	"strings"

	"github.com/kelseyhightower/envconfig"

	"github.com/kiwiworks/rodent/system/errors"
)

func FromEnv[T any](prefixes ...string) (T, error) {
	var t T
	if err := envconfig.Process(strings.Join(prefixes, "_"), &t); err != nil {
		return t, errors.Wrapf(err, "failed to load config from environment")
	}
	return t, nil
}
