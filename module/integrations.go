package module

import (
	"github.com/pkg/errors"
	"go.uber.org/fx"

	"github.com/kiwiworks/rodent/assert"
	"github.com/kiwiworks/rodent/command"
	"github.com/kiwiworks/rodent/database/migration"
	"github.com/kiwiworks/rodent/slices"
	"github.com/kiwiworks/rodent/system/opt"
	"github.com/kiwiworks/rodent/web/api"
)

func annotateAndAppend(groupTag string, factories []any, opt *Module) {
	annotatedFactories := slices.Map(factories, func(in any) any {
		return fx.Annotate(in, fx.ResultTags(groupTag))
	})
	opt.Public = append(opt.Public, annotatedFactories...)
}

func Handlers(handlerProviders ...any) opt.Option[Module] {
	for _, handlerProvider := range handlerProviders {
		if err := assert.FuncHasReturn[*api.Handler](handlerProvider); err != nil {
			panic(errors.Wrap(err, "module.Handlers only accepts function of any arity, which must only return *api.Handler"))
		}
	}
	return func(opt *Module) {
		annotateAndAppend(`group:"api.handler"`, handlerProviders, opt)
	}
}

func Migrations(migrationProviders ...any) opt.Option[Module] {
	for _, migrationProvider := range migrationProviders {
		if err := assert.FuncHasReturn[*migration.Migration](migrationProvider); err != nil {
			panic(errors.Wrap(err, "module.Migrations only accepts function of any arity, which must only return *migration.Migration"))
		}
	}
	return func(opt *Module) {
		annotateAndAppend(`group:"migration.migration"`, migrationProviders, opt)
	}
}

func Commands(commandProviders ...any) opt.Option[Module] {
	for _, commandProvider := range commandProviders {
		if err := assert.FuncHasReturn[*command.Command](commandProvider); err != nil {
			panic(errors.Wrap(err, "module.Commands only accepts function of any arity, which must only return *command.Command"))
		}
	}
	return func(opt *Module) {
		annotateAndAppend(`group:"command"`, commandProviders, opt)
	}
}
