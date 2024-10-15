package auth

import (
	"github.com/kiwiworks/rodent/app"
	"github.com/kiwiworks/rodent/app/module"
)

func Module() app.Module {
	return app.NewModule(
		module.Public(
			NewMiddleware,
			NewProviders,
		),
	)
}
