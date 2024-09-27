package auth

import "github.com/kiwiworks/rodent/module"

func Module() module.Module {
	return module.New(
		"rodent.auth",
		module.Public(NewMiddleware),
	)
}
