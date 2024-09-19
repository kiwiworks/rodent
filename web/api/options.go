package api

import "github.com/kiwiworks/rodent/system/opt"

func Protected() opt.Option[Options] {
	return func(opt *Options) {
		opt.Protected = true
	}
}

func Tags(tags ...string) opt.Option[Options] {
	return func(opt *Options) {
		opt.Tags = tags
	}
}
