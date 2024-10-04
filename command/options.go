package command

import (
	"github.com/spf13/cobra"

	"github.com/kiwiworks/rodent/errors"
	"github.com/kiwiworks/rodent/maps"
	"github.com/kiwiworks/rodent/system/opt"
)

func Runner(impl func(cmd *cobra.Command, args []string) error) opt.Option[Command] {
	return func(opt *Command) {
		opt.Run = impl
	}
}

type Flag struct {
	Name        string
	Shorthand   string
	Usage       string
	OneRequired bool
	Required    bool
	Exclusive   bool
}

func flagOpt(flag Flag, handler func(cmd *cobra.Command) error) opt.Option[Command] {
	return func(opt *Command) {
		if flag.Exclusive {
			opt.MutuallyExclusiveFlags = append(opt.MutuallyExclusiveFlags, flag.Name)
		}
		if flag.OneRequired {
			opt.RequiredOneFlags = append(opt.RequiredOneFlags, flag.Name)
		}
		if flag.Required {
			opt.RequiredFlags = append(opt.RequiredFlags, flag.Name)
		}
		opt.FlagHandlers[flag.Name] = handler
	}
}

func StringFlag(flag Flag, ptr *string) opt.Option[Command] {
	return flagOpt(flag, func(cmd *cobra.Command) error {
		if ptr == nil {
			return errors.Newf("command.StringFlag expects a pointer to a string, but got nil")
		}
		cmd.Flags().StringVarP(ptr, flag.Name, flag.Shorthand, *ptr, flag.Usage)
		return nil
	})
}

func BoolFlag(flag Flag, ptr *bool) opt.Option[Command] {
	return flagOpt(flag, func(cmd *cobra.Command) error {
		if ptr == nil {
			return errors.Newf("command.BoolFlag expects a pointer to a bool, but got nil")
		}
		cmd.Flags().BoolVarP(ptr, flag.Name, flag.Shorthand, *ptr, flag.Usage)
		return nil
	})
}

func IntFlag(flag Flag, ptr *int) opt.Option[Command] {
	return flagOpt(flag, func(cmd *cobra.Command) error {
		if ptr == nil {
			return errors.Newf("command.IntFlag expects a pointer to an int, but got nil")
		}
		cmd.Flags().IntVarP(ptr, flag.Name, flag.Shorthand, *ptr, flag.Usage)
		return nil
	})
}

func Example(example string) opt.Option[Command] {
	return func(opt *Command) {
		opt.Example = example
	}
}

func Annotation(key, value string) opt.Option[Command] {
	return func(opt *Command) {
		opt.Annotations[key] = value
	}
}

func Annotations(annotations map[string]string) opt.Option[Command] {
	return func(opt *Command) {
		opt.Annotations = maps.Merged(opt.Annotations, annotations)
	}
}

func Deprecated(reason ...string) opt.Option[Command] {
	return func(opt *Command) {
		if len(reason) > 0 {
			opt.Deprecated = reason[0]
		}
		opt.Deprecated = "This command is deprecated."
	}
}

func SuggestFor(alternatives ...string) opt.Option[Command] {
	return func(opt *Command) {
		opt.SuggestFor = alternatives
	}
}
