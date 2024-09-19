package api

import "rodent/system/opt"

type Config struct {
	ErrorConverter func(err error) error
}

func ErrorConverter(impl func(err error) error) opt.Option[Config] {
	return func(opt *Config) {
		opt.ErrorConverter = impl
	}
}

func NewConfig(opts ...opt.Option[Config]) *Config {
	cfg := DefaultConfig()
	opt.Apply(cfg, opts...)
	return cfg
}

func NoOpErrorConverter(err error) error {
	return err
}

func DefaultConfig() *Config {
	return &Config{
		ErrorConverter: NoOpErrorConverter,
	}
}
