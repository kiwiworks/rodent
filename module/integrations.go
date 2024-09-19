package module

import (
	"go.uber.org/fx"

	"rodent/slices"
	"rodent/system/opt"
)

func annotateAndAppend(groupTag string, factories []any, opt *Module) {
	annotatedFactories := slices.Map(factories, func(in any) any {
		return fx.Annotate(in, fx.ResultTags(groupTag))
	})
	opt.Public = append(opt.Public, annotatedFactories...)
}

func Handlers(factories ...any) opt.Option[Module] {
	return func(opt *Module) {
		annotateAndAppend(`group:"api.handler"`, factories, opt)
	}
}

func Migrations(factories ...any) opt.Option[Module] {
	return func(opt *Module) {
		annotateAndAppend(`group:"migration.migration"`, factories, opt)
	}
}

func Commands(commands ...any) opt.Option[Module] {
	return func(opt *Module) {
		annotateAndAppend(`group:"command"`, commands, opt)
	}
}
