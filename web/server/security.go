package server

import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"

	"rodent/system/opt"
)

type SecurityConfig struct {
	Name           string
	Scheme         *huma.SecurityScheme
	InjectedValues map[string]reflect.Type
}

func InjectValue[T any](name string) opt.Option[SecurityConfig] {
	return func(opt *SecurityConfig) {
		var v T
		t := reflect.TypeOf(v)
		opt.InjectedValues[name] = t
	}
}

func NewSecurityConfig(name string, scheme *huma.SecurityScheme, opts ...opt.Option[SecurityConfig]) *SecurityConfig {
	cfg := &SecurityConfig{
		Name:           name,
		Scheme:         scheme,
		InjectedValues: make(map[string]reflect.Type),
	}
	opt.Apply(cfg, opts...)
	return cfg
}
