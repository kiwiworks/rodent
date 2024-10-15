package web

import (
	"github.com/kiwiworks/rodent/app"
	"github.com/kiwiworks/rodent/app/module"
	"github.com/kiwiworks/rodent/web/auth"
	"github.com/kiwiworks/rodent/web/server"
)

func Module() app.Module {
	return app.NewModule(
		module.Public(
			server.NewMux,
			server.NewHuma,
			server.NewRouter,
			server.New,
		),
		module.Service[server.Server](),
		module.SubModules(auth.Module),
	)
}
