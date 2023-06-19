package users

import (
	"context"
	"database/sql"
	"fmt"

	_ "embed"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Service struct {
	postgresURL   string
	migrationsURL string
	*sql.DB
	stmt_select_user_by_uuid  *sql.Stmt
	stmt_select_user_by_email *sql.Stmt
	stmt_delete_user_by_uuid  *sql.Stmt
}

func NewService(postgresURL, migrationsURL string) *Service {
	return &Service{
		postgresURL:   postgresURL,
		migrationsURL: migrationsURL,
	}
}

func (s *Service) Connect(ctx context.Context) error {
	db, err := sql.Open("postgres", s.postgresURL)
	if err != nil {
		return fmt.Errorf("failed to open db connection: %w", err)
	}
	s.DB = db

	err = s.migrateSqlDatabase()
	if err != nil {
		return fmt.Errorf("user db migration failed: %w", err)
	}

	return s.prepareSqlStatements(ctx)
}

// Close all open connections safely.
func (s *Service) Close() (err error) {
	err = s.DB.Close()
	if err != nil {
		return err
	}
	return
}

// migrateSqlDatabase runs sql schema statements in the migrations directory.
// db schema must be immutable.
// sql file names must take the form of: <index>_<title>.<up|down>.sql
//
// sql files will be run in the order of their index.
func (s *Service) migrateSqlDatabase() (err error) {
	driver, err := postgres.WithInstance(s.DB, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(s.migrationsURL, "postgres", driver)
	if err != nil {
		return err
	}
	m.Up()
	return
}

func (s *Service) prepareSqlStatements(ctx context.Context) error {
	err := s.prepareGetByUUID(ctx)
	if err != nil {
		return err
	}
	err = s.prepareGetByEmail(ctx)
	if err != nil {
		return err
	}
	err = s.prepareDeleteByUUID(ctx)
	if err != nil {
		return err
	}
	return err
}
