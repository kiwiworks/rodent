package command

import (
	"github.com/spf13/cobra"
)

type (
	Command struct {
		Name  string
		Short string
		Long  string
	}
)

func New(name string, short string, long string) *Command {
	return &Command{
		Name:  name,
		Short: short,
		Long:  long,
	}
}

func (c *Command) asCobraCommand() *cobra.Command {
	return &cobra.Command{
		Use:        "",
		Aliases:    nil,
		SuggestFor: nil,
		Short:      c.Short,
		GroupID:    "",
		Long:       c.Long,
		Example:    "",
		ValidArgs:  []string{},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) (
			[]string,
			cobra.ShellCompDirective,
		) {
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		Args:                   func(cmd *cobra.Command, args []string) error { return nil },
		ArgAliases:             []string{},
		BashCompletionFunction: "",
		//Deprecated:                 "",
		Annotations:       nil,
		Version:           "",
		PersistentPreRun:  nil,
		PersistentPreRunE: nil,
		PreRun:            nil,
		PreRunE:           nil,
		Run: func(cmd *cobra.Command, args []string) {

		},
		RunE:                       nil,
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
}
