package app

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	"github.com/kiwiworks/rodent/module"
	"github.com/kiwiworks/rodent/system/logger"
	"github.com/kiwiworks/rodent/system/manifest"
	"github.com/kiwiworks/rodent/system/opt"
)

type App struct {
	manifest  *manifest.Manifest
	fxOptions []fx.Option
	di        *fx.App
}

func Modules(modules ...func() module.Module) opt.Option[App] {
	return func(opt *App) {
		for _, m := range modules {
			opt.fxOptions = append(opt.fxOptions, m().IntoFxModule())
		}
	}
}

func fxLogProvider() fxevent.Logger {
	log := &fxevent.ZapLogger{
		Logger: logger.New(),
	}
	log.UseLogLevel(log.Logger.Level())
	return log
}

func New(name, version string, opts ...opt.Option[App]) *App {
	m := manifest.New(name, version)
	app := &App{
		manifest: m,
		fxOptions: []fx.Option{
			fx.Supply(m),
			fx.WithLogger(fxLogProvider),
		},
	}
	opt.Apply(app, opts...)
	app.di = fx.New(app.fxOptions...)
	return app
}

func (app *App) Run() {
	ctx := context.Background()
	log := logger.New()
	log.Info("starting application")
	if err := app.di.Start(ctx); err != nil {
		panic(err)
	}
	signal := <-app.di.Wait()
	log.Info("application stopped", zap.Int("exitCode", signal.ExitCode), zap.Stringer("signal", signal.Signal))
}
