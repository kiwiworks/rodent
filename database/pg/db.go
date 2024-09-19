package pg

import (
	"database/sql"

	entsql "entgo.io/ent/dialect/sql"
)

const Dialect = "postgres"

type Database struct {
	db *sql.DB
}

func NewDatabase(db *sql.DB) *Database {
	return &Database{db: db}
}

func (d Database) Driver() *entsql.Driver {
	return entsql.OpenDB(Dialect, d.db)
}
