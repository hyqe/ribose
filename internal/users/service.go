package users

import (
	"context"
	"database/sql"
	"fmt"

	_ "embed"

	"github.com/go-playground/validator/v10"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type UserStorage struct {
	postgresURL   string
	migrationsURL string
	*validator.Validate
	*sql.DB
	stmt_select_user_by_uuid  *sql.Stmt
	stmt_select_user_by_email *sql.Stmt
	stmt_delete_user_by_uuid  *sql.Stmt
}

func NewUserStorage(postgresURL, migrationsURL string) *UserStorage {
	return &UserStorage{
		postgresURL:   postgresURL,
		migrationsURL: migrationsURL,
		Validate:      validator.New(),
	}
}

func (u *UserStorage) Connect(ctx context.Context) error {
	db, err := sql.Open("postgres", u.postgresURL)
	if err != nil {
		return fmt.Errorf("failed to open db connection: %w", err)
	}
	u.DB = db

	err = u.migrateSqlDatabase()
	if err != nil {
		return fmt.Errorf("user db migration failed: %w", err)
	}

	return u.prepareSqlStatements(ctx)
}

// Close all open connections safely.
func (u *UserStorage) Close() (err error) {
	err = u.DB.Close()
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
func (u *UserStorage) migrateSqlDatabase() (err error) {
	driver, err := postgres.WithInstance(u.DB, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(u.migrationsURL, "postgres", driver)
	if err != nil {
		return err
	}
	m.Up()
	return
}

func (u *UserStorage) prepareSqlStatements(ctx context.Context) error {
	err := u.prepareGetByUUID(ctx)
	if err != nil {
		return err
	}
	err = u.prepareGetByEmail(ctx)
	if err != nil {
		return err
	}
	err = u.prepareDeleteByUUID(ctx)
	if err != nil {
		return err
	}
	return err
}
