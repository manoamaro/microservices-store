package infra

import (
	"database/sql"
	"embed"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"io/fs"
	"log"
)

type Migrator interface {
	Up() error
	Down() error
}

type migrator struct {
	migration *migrate.Migrate
	fs        embed.FS
}

func NewMigrator(dbUrl string, fs fs.FS) Migrator {
	if db, err := sql.Open("postgres", dbUrl); err != nil {
		log.Fatal(err)
		return nil
	} else if migrationsSource, err := iofs.New(fs, "migrations"); err != nil {
		log.Fatal(err)
		return nil
	} else if driver, err := postgres.WithInstance(db, &postgres.Config{}); err != nil {
		log.Fatal(err)
		return nil
	} else if migrations, err := migrate.NewWithInstance("iofs", migrationsSource, "postgres", driver); err != nil {
		log.Fatal(err)
		return nil
	} else {
		return &migrator{migration: migrations}
	}
}

func (m *migrator) Up() error {
	return m.migration.Up()
}

func (m *migrator) Down() error {
	return m.migration.Down()
}
