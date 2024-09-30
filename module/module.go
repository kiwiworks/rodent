package module

import (
	"go.uber.org/fx"

	"github.com/kiwiworks/rodent/slices"
	"github.com/kiwiworks/rodent/system/opt"
)

type (
	IModule interface {
		IntoFxModule() fx.Option
	}
	Module struct {
		Name       string
		Public     []any
		Private    []any
		Instances  []any
		Decorators []any
		Invokers   []any
		SubModules []IModule
	}
)

func New(name string, opts ...opt.Option[Module]) Module {
	mod := Module{
		Name:       name,
		Public:     []any{},
		Private:    []any{},
		Instances:  []any{},
		Decorators: []any{},
		Invokers:   []any{},
		SubModules: []IModule{},
	}
	opt.Apply(&mod, opts...)
	return mod
}

func (m Module) IntoFxModule() fx.Option {
	opts := make([]fx.Option, 0)
	opts = append(opts, fx.Provide(m.Public...))
	opts = append(opts, fx.Provide(append(m.Private, fx.Private)...))
	opts = append(opts, fx.Decorate(m.Decorators...))
	opts = append(opts, fx.Invoke(m.Invokers...))
	opts = append(opts, fx.Supply(m.Instances...))
	opts = append(opts, slices.Map(m.SubModules, func(in IModule) fx.Option {
		return in.IntoFxModule()
	})...)

	return fx.Module(
		m.Name,
		opts...,
	)
}
