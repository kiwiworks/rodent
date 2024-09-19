package opt

type (
	Option[Options any]    func(opt *Options)
	TryOption[Options any] func(opt *Options) error
)

func Apply[Options any](this *Options, opts ...Option[Options]) {
	for _, opt := range opts {
		opt(this)
	}
}

func TryApply[Options any](this *Options, opts ...TryOption[Options]) error {
	for _, opt := range opts {
		if err := opt(this); err != nil {
			return err
		}
	}
	return nil
}
