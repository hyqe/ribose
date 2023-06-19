package server

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func connectDB(ctx context.Context, postgresURL, migrationsURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", postgresURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open db connection: %w", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}
	m, err := migrate.NewWithDatabaseInstance(migrationsURL, "postgres", driver)
	if err != nil {
		return nil, err
	}
	m.Up()

	return db, nil
}
