package migration

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"log"
)

type MigrateData struct {
	m *migrate.Migrate
}

func New(db *sql.DB, dirPath string) (*MigrateData, error) {
	var m MigrateData

	driver, err := postgres.WithInstance(db, &postgres.Config{})

	if err != nil {
		log.Fatal("Failed to get driver")
	}

	m.m, err = migrate.NewWithDatabaseInstance(dirPath, "postgres", driver)

	return &m, err
}

func (m *MigrateData) Up() error {
	err := m.m.Up()
	if err != nil {
		return err
	}
	return nil
}

func (m *MigrateData) Down() error {
	err := m.m.Down()
	if err != nil {
		return err
	}
	return nil
}

func (m *MigrateData) Force(version int) error {
	err := m.m.Force(version)
	if err != nil {
		return err
	}
	return nil
}
