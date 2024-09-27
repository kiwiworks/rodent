package config

import (
	"strings"

	"github.com/kelseyhightower/envconfig"

	"github.com/kiwiworks/rodent/errors"
)

func FromEnv[T any](prefixes ...string) (T, error) {
	var t T
	key := strings.ToUpper(strings.Join(prefixes, "_"))
	if err := envconfig.Process(key, &t); err != nil {
		return t, errors.Wrapf(err, "failed to load config from environment variable %s", key)
	}
	return t, nil
}
