package internal

import (
	"database/sql"
	"embed"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var fs embed.FS

type Migration struct {
	migration *migrate.Migrate
}

func NewMigration(db *sql.DB) (*Migration, error) {
	migrationsSource, err := iofs.New(fs, "migrations")
	if err != nil {
		return nil, err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	migrations, err := migrate.NewWithInstance("iofs", migrationsSource, "postgres", driver)
	if err != nil {
		return nil, err
	}

	m := &Migration{
		migration: migrations,
	}

	return m, nil
}

func (m *Migration) Up() error {
	return m.migration.Up()
}
