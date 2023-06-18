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
	postgresDSN   string
	migrationsURL string
	*sql.DB
	stmt_select_user_by_uuid  *sql.Stmt
	stmt_select_user_by_email *sql.Stmt
	stmt_delete_user_by_uuid  *sql.Stmt
}

func NewService(postgresDSN, migrationsURL string) *Service {
	return &Service{
		postgresDSN:   postgresDSN,
		migrationsURL: migrationsURL,
	}
}

func (s *Service) Connect(ctx context.Context) error {
	db, err := sql.Open("postgres", s.postgresDSN)
	if err != nil {
		return fmt.Errorf("failed to open db connection: %w", err)
	}
	s.DB = db

	err = s.migrate()
	if err != nil {
		return fmt.Errorf("user db migration failed: %w", err)
	}

	return s.prepare(ctx)
}

// Close all open connections safely.
func (s *Service) Close() (err error) {
	err = s.DB.Close()
	if err != nil {
		return err
	}
	return
}

func (s *Service) migrate() (err error) {
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

func (s *Service) prepare(ctx context.Context) error {
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
