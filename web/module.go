package web

import (
	"github.com/kiwiworks/rodent/module"
	"github.com/kiwiworks/rodent/web/server"
)

func Module() module.Module {
	return module.New(
		"core.web",
		module.Public(
			server.NewMux,
			server.NewHuma,
			server.NewRouter,
			server.New,
		),
		module.Service[server.Server](),
	)
}
