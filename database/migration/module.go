package migration

import "rodent/module"

func Module() module.Module {
	return module.New(
		"database.migration",
		module.Public(NewMigrator),
		module.Service[Migrator](),
	)
}
