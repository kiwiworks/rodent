package command

import (
	"context"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/kiwiworks/rodent/logger"
	"github.com/kiwiworks/rodent/system/manifest"
)

type Root struct {
	shutdowner fx.Shutdowner
	cancel     context.CancelFunc
	root       *cobra.Command
}

type RootParams struct {
	fx.In
	Manifest *manifest.Manifest
	Shutdown fx.Shutdowner
	Commands []*Command `group:"command"`
}

func NewRoot(params RootParams) *Root {
	log := logger.New()

	rootCmd := &cobra.Command{
		Use:     params.Manifest.Application,
		Version: params.Manifest.Version.String(),
	}
	for _, cmd := range params.Commands {
		cobraCmd := cmd.asCobraCommand()
		log.Debug("added new command", zap.String("command.short", cobraCmd.Short))
		rootCmd.AddCommand(cobraCmd)
	}
	return &Root{
		root:       rootCmd,
		shutdowner: params.Shutdown,
	}
}

func (r *Root) OnStart(ctx context.Context) error {
	log := logger.FromContext(ctx)
	ctx, r.cancel = context.WithCancel(context.Background())
	go func() {
		defer r.cancel()
		if err := r.root.ExecuteContext(ctx); err != nil {
			log.Error("Execution failed", zap.Error(err))
			if err = r.shutdowner.Shutdown(fx.ExitCode(1)); err != nil {
				log.Fatal("Failed to shutdown application", zap.Error(err))
			}
		}
		if err := r.shutdowner.Shutdown(); err != nil {
			log.Fatal("Failed to shutdown application", zap.Error(err))
		}
	}()
	return nil
}

func (r *Root) OnStop(context.Context) error {
	if r.cancel != nil {
		r.cancel()
	}
	return nil
}
