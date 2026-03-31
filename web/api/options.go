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

func Description(description string) opt.Option[Options] {
	return func(opt *Options) {
		opt.Description = description
	}
}

func OperationID(operationID string) opt.Option[Options] {
	return func(opt *Options) {
		opt.OperationId = operationID
	}
}

// Deprecated marks this endpoint as deprecated in the OpenAPI spec.
// successorPath is the v2 equivalent (e.g., "/api/v2/persons") — included in
// the Deprecation response header as a Link rel="successor-version".
func Deprecated(successorPath string) opt.Option[Options] {
	return func(opt *Options) {
		opt.Deprecated = true
		opt.SuccessorVersion = successorPath
	}
}
