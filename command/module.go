package command

import "rodent/module"

func Module() module.Module {
	return module.New(
		"core.command",
		module.Public(NewRoot),
		module.Service[Root](),
	)
}
