package migrations

import (
	"embed"

	"github.com/uptrace/bun/migrate"
)

var Migrations = migrate.NewMigrations()

//go:embed *.sql
var sqlMigrations embed.FS

//nolint:gochecknoinits // Migrations should be initialized in init functions.
func init() {
	if err := Migrations.Discover(sqlMigrations); err != nil {
		panic(err)
	}
}
