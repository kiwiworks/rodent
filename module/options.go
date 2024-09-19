package module

import (
	"github.com/kiwiworks/rodent/slices"
	"github.com/kiwiworks/rodent/system/opt"
)

func Public(providers ...any) opt.Option[Module] {
	return func(opt *Module) {
		opt.Public = append(opt.Public, providers...)
	}
}

func Private(providers ...any) opt.Option[Module] {
	return func(opt *Module) {
		opt.Private = append(opt.Private, providers...)
	}
}

func Supply(suppliers ...any) opt.Option[Module] {
	return func(opt *Module) {
		opt.Instances = append(opt.Instances, suppliers...)
	}
}

func SubModules(modules ...func() Module) opt.Option[Module] {
	return func(opt *Module) {
		opt.SubModules = append(opt.SubModules, slices.Map(modules, func(in func() Module) Module {
			return in()
		})...)
	}
}
