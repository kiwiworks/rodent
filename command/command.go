package command

import (
	"github.com/spf13/cobra"

	"github.com/kiwiworks/rodent/errors"
	"github.com/kiwiworks/rodent/system/opt"
)

type (
	Command struct {
		Name                   string
		Short                  string
		Long                   string
		Run                    func(cmd *cobra.Command, args []string) error
		Example                string
		Annotations            map[string]string
		Deprecated             string
		SuggestFor             []string
		MutuallyExclusiveFlags []string
		RequiredOneFlags       []string
		RequiredFlags          []string
		FlagHandlers           map[string]func(cmd *cobra.Command) error
	}
)

func New(name string, short string, long string, opts ...opt.Option[Command]) *Command {
	cmd := &Command{
		Name:  name,
		Short: short,
		Long:  long,
		Run: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Annotations:            map[string]string{},
		SuggestFor:             []string{},
		MutuallyExclusiveFlags: []string{},
		RequiredOneFlags:       []string{},
		RequiredFlags:          []string{},
		FlagHandlers:           map[string]func(cmd *cobra.Command) error{},
	}
	opt.Apply(cmd, opts...)
	return cmd
}

func (c *Command) asCobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:        c.Name,
		Aliases:    nil,
		SuggestFor: c.SuggestFor,
		Short:      c.Short,
		GroupID:    "",
		Long:       c.Long,
		Example:    c.Example,
		//ValidArgs:  []string{},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) (
			[]string,
			cobra.ShellCompDirective,
		) {
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		Args:                   func(cmd *cobra.Command, args []string) error { return nil },
		ArgAliases:             []string{},
		BashCompletionFunction: "",
		Deprecated:             c.Deprecated,
		Annotations:            c.Annotations,
		Version:                "",
		PersistentPreRun:       nil,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			cmd.MarkFlagsOneRequired(c.RequiredOneFlags...)
			cmd.MarkFlagsMutuallyExclusive(c.MutuallyExclusiveFlags...)
			for _, flag := range c.RequiredFlags {
				if err := cmd.MarkFlagRequired(flag); err != nil {
					return errors.Wrapf(err, "failed to mark flag '%s' as required", flag)
				}
			}

			return nil
		},
		PreRun:  nil,
		PreRunE: nil,
		Run:     nil,
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Run(cmd, args)
		},
		PostRun:                    nil,
		PostRunE:                   nil,
		PersistentPostRun:          nil,
		PersistentPostRunE:         nil,
		FParseErrWhitelist:         cobra.FParseErrWhitelist{},
		CompletionOptions:          cobra.CompletionOptions{},
		TraverseChildren:           false,
		Hidden:                     false,
		SilenceErrors:              false,
		SilenceUsage:               false,
		DisableFlagParsing:         false,
		DisableAutoGenTag:          false,
		DisableFlagsInUseLine:      false,
		DisableSuggestions:         false,
		SuggestionsMinimumDistance: 0,
	}
	for name, handler := range c.FlagHandlers {
		if err := handler(cmd); err != nil {
			panic(errors.Wrapf(err, "failed to handle flag parsing for '%s'", name))
		}
	}

	return cmd
}
