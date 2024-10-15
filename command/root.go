package command

import (
	"context"
	"strings"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/kiwiworks/rodent/errors"
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

func NewRoot(params RootParams) (*Root, error) {
	log := logger.New()

	rootCmd := &cobra.Command{
		Use:     params.Manifest.Application,
		Version: params.Manifest.Version.String(),
	}
	allCommands := make(map[string]*Command)
	allCobraCommands := make(map[string]*cobra.Command)
	edges := make(map[string][]string)
	allCobraCommands[rootCmd.Name()] = rootCmd
	for _, cmd := range params.Commands {
		allCommands[cmd.Name] = cmd
		if cmd.ChildOf != nil && len(cmd.ChildOf) > 0 {
			parentCount := len(cmd.ChildOf)
			for idx, parent := range cmd.ChildOf {
				var child string
				if idx == parentCount-1 {
					child = cmd.Name
				} else {
					child = cmd.ChildOf[idx+1]
				}
				edges[parent] = append(edges[parent], child)
			}
		} else {
			edges[rootCmd.Name()] = append(edges[rootCmd.Name()], cmd.Name)
		}
		cobraCmd := cmd.asCobraCommand()
		allCobraCommands[cmd.Name] = cobraCmd
		log.Debug("added new command", zap.String("command.short", cobraCmd.Short))
	}
	for parent, children := range edges {
		parentCmd, exists := allCobraCommands[parent]
		if !exists {
			return nil, errors.Newf(
				"parent command %s does not exist, but is needed by the following children `%s`",
				parent, strings.Join(children, ", "),
			)
		}
		for _, child := range children {
			childCmd, exists := allCobraCommands[child]
			if !exists {
				return nil, errors.Newf(
					"child command %s does not exist, but was declared as parented to %s",
					child, parent,
				)
			}
			parentCmd.AddCommand(childCmd)
		}
	}

	return &Root{
		root:       rootCmd,
		shutdowner: params.Shutdown,
	}, nil
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
