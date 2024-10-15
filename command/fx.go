package command

import (
	"github.com/pkg/errors"
	"go.uber.org/fx"

	"github.com/kiwiworks/rodent/app"
	"github.com/kiwiworks/rodent/assert"
	"github.com/kiwiworks/rodent/slices"
	"github.com/kiwiworks/rodent/system/opt"
)

func annotateAndAppend(groupTag string, factories []any, opt *app.Module) {
	annotatedFactories := slices.Map(factories, func(in any) any {
		return fx.Annotate(in, fx.ResultTags(groupTag))
	})
	opt.Public = append(opt.Public, annotatedFactories...)
}

func Commands(commandProviders ...any) opt.Option[app.Module] {
	for _, commandProvider := range commandProviders {
		if err := assert.FuncHasReturn[*Command](commandProvider); err != nil {
			panic(errors.Wrap(err, "module.Commands only accepts function of any arity, which must only return *command.Command"))
		}
	}
	return func(opt *app.Module) {
		annotateAndAppend(`group:"command"`, commandProviders, opt)
	}
}
