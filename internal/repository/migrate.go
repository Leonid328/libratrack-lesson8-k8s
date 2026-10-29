package repository

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(databaseURL, path string) error {
	if path == "" {
		path = "migrations"
	}
	m, err := migrate.New("file://"+path, databaseURL)
	if err != nil {
		return fmt.Errorf("migrate open: %w", err)
	}
	defer m.Close()
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
