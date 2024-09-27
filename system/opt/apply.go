package opt

// Option represents a function that configures an instance of the specified type Options by modifying its fields.
type Option[Options any] func(opt *Options)

// TryOption defines a function type that modifies an Options object and returns an error if the modification fails.
type TryOption[Options any] func(opt *Options) error

// Apply sets the given options on the provided structure.
// this is a pointer to the structure to configure.
// opts are a list of configuration options to apply.
func Apply[Options any](this *Options, opts ...Option[Options]) {
	for _, opt := range opts {
		opt(this)
	}
}

// TryApply applies a list of TryOption functions to the given options and returns an error if any of the functions fail.
func TryApply[Options any](this *Options, opts ...TryOption[Options]) error {
	for _, opt := range opts {
		if err := opt(this); err != nil {
			return err
		}
	}
	return nil
}
