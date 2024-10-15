package command

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/kiwiworks/rodent/errors"
	"github.com/kiwiworks/rodent/maps"
	"github.com/kiwiworks/rodent/system/opt"
)

// Runner is where your actual command code should go.
// Deprecated: use command.Do instead
func Runner(impl func(cmd *cobra.Command, args []string) error) opt.Option[Command] {
	return func(opt *Command) {
		opt.Run = impl
	}
}

// Do is where your actual command code should go.
func Do(impl func(ctx context.Context) error) opt.Option[Command] {
	return func(opt *Command) {
		opt.Run = func(cmd *cobra.Command, args []string) error {
			return impl(cmd.Context())
		}
	}
}

// Flag represents a command-line option for a CLI application.
type Flag struct {
	// Name is the identifier for the command-line option.
	Name string
	// Shorthand is a single-character alias for the flag that allows for quicker command-line access.
	Shorthand string
	// Usage describes the purpose or functionality of the command-line flag.
	Usage string
	// OneRequired indicates that at least one of the flags marked with this field must be specified.
	OneRequired bool
	// Required indicates whether this flag is mandatory for the command-line tool to execute successfully.
	Required bool
	// Exclusive indicates whether the flag is mutually exclusive with other flags in the command-line tool.
	Exclusive bool
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

// StringFlag configures a string flag for a command using the provided Flag configuration and a pointer to store the flag's value.
func StringFlag(flag Flag, ptr *string) opt.Option[Command] {
	return flagOpt(flag, func(cmd *cobra.Command) error {
		if ptr == nil {
			return errors.Newf("command.StringFlag expects a pointer to a string, but got nil")
		}
		cmd.Flags().StringVarP(ptr, flag.Name, flag.Shorthand, *ptr, flag.Usage)
		return nil
	})
}

// StringsFlag registers a string slice command-line flag for the given Command.
func StringsFlag(flag Flag, ptr *[]string) opt.Option[Command] {
	return flagOpt(flag, func(cmd *cobra.Command) error {
		if ptr == nil {
			return errors.Newf("command.StringsFlag expects a pointer to a string slice, but got nil")
		}
		cmd.Flags().StringSliceVarP(ptr, flag.Name, flag.Shorthand, *ptr, flag.Usage)
		return nil
	})
}

// BoolFlag creates an option for a boolean flag within a Cobra command and links it to a specified bool pointer.
func BoolFlag(flag Flag, ptr *bool) opt.Option[Command] {
	return flagOpt(flag, func(cmd *cobra.Command) error {
		if ptr == nil {
			return errors.Newf("command.BoolFlag expects a pointer to a bool, but got nil")
		}
		cmd.Flags().BoolVarP(ptr, flag.Name, flag.Shorthand, *ptr, flag.Usage)
		return nil
	})
}

// IntFlag registers an integer flag with the given flag specifications and stores the value in the specified pointer.
func IntFlag(flag Flag, ptr *int) opt.Option[Command] {
	return flagOpt(flag, func(cmd *cobra.Command) error {
		if ptr == nil {
			return errors.Newf("command.IntFlag expects a pointer to an int, but got nil")
		}
		cmd.Flags().IntVarP(ptr, flag.Name, flag.Shorthand, *ptr, flag.Usage)
		return nil
	})
}

// Example sets the example string for a Command instance.
func Example(example string) opt.Option[Command] {
	return func(opt *Command) {
		opt.Example = example
	}
}

// Annotation sets a key-value pair in the Annotations map of a Command.
func Annotation(key, value string) opt.Option[Command] {
	return func(opt *Command) {
		opt.Annotations[key] = value
	}
}

// Annotations sets additional metadata for the command via the provided annotations map.
func Annotations(annotations map[string]string) opt.Option[Command] {
	return func(opt *Command) {
		opt.Annotations = maps.Merged(opt.Annotations, annotations)
	}
}

// Deprecated marks a Command as deprecated optionally specifying a reason.
func Deprecated(reason ...string) opt.Option[Command] {
	return func(opt *Command) {
		if len(reason) > 0 {
			opt.Deprecated = reason[0]
		}
		opt.Deprecated = "This command is deprecated."
	}
}

// SuggestFor returns an Option that sets the SuggestFor field of a Command to the provided alternatives.
func SuggestFor(alternatives ...string) opt.Option[Command] {
	return func(opt *Command) {
		opt.SuggestFor = alternatives
	}
}
