package api

import "github.com/kiwiworks/rodent/system/opt"

func Auth(providerNames ...string) opt.Option[Options] {
	return func(opt *Options) {
		opt.AuthProviders = providerNames
	}
}

func Oauth2(providerName string, scopes ...string) opt.Option[Options] {
	return func(opt *Options) {
		opt.OAuth2Providers[providerName] = scopes
	}
}

func Tags(tags ...string) opt.Option[Options] {
	return func(opt *Options) {
		opt.Tags = tags
	}
}
