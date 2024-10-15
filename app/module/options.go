package module

import (
	"github.com/kiwiworks/rodent/app"
	"github.com/kiwiworks/rodent/slices"
	"github.com/kiwiworks/rodent/system/opt"
)

func Public(providers ...any) opt.Option[app.Module] {
	return func(opt *app.Module) {
		opt.Public = append(opt.Public, providers...)
	}
}

func Private(providers ...any) opt.Option[app.Module] {
	return func(opt *app.Module) {
		opt.Private = append(opt.Private, providers...)
	}
}

func Supply(suppliers ...any) opt.Option[app.Module] {
	return func(opt *app.Module) {
		opt.Instances = append(opt.Instances, suppliers...)
	}
}

func SubModules(modules ...func() app.Module) opt.Option[app.Module] {
	return func(opt *app.Module) {
		opt.SubModules = append(opt.SubModules, slices.Map(modules, func(in func() app.Module) app.IModule {
			return in()
		})...)
	}
}
