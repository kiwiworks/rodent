package module

import (
	"github.com/pkg/errors"
	"go.uber.org/fx"

	"github.com/kiwiworks/rodent/app"
	"github.com/kiwiworks/rodent/assert"
	"github.com/kiwiworks/rodent/slices"
	"github.com/kiwiworks/rodent/system/opt"
	"github.com/kiwiworks/rodent/web/api"
)

func annotateAndAppend(groupTag string, factories []any, opt *app.Module) {
	annotatedFactories := slices.Map(factories, func(in any) any {
		return fx.Annotate(in, fx.ResultTags(groupTag))
	})
	opt.Public = append(opt.Public, annotatedFactories...)
}

func Handlers(handlerProviders ...any) opt.Option[app.Module] {
	for _, handlerProvider := range handlerProviders {
		if err := assert.FuncHasReturn[*api.Handler](handlerProvider); err != nil {
			panic(errors.Wrap(err, "module.Handlers only accepts function of any arity, which must only return *api.Handler"))
		}
	}
	return func(opt *app.Module) {
		annotateAndAppend(`group:"api.handler"`, handlerProviders, opt)
	}
}
