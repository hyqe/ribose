package users

import (
	"context"
	"database/sql"

	"github.com/go-playground/validator/v10"
)

type UserStorage struct {
	*validator.Validate
	*sql.DB
	stmt_select_user_by_uuid  *sql.Stmt
	stmt_select_user_by_email *sql.Stmt
	stmt_delete_user_by_uuid  *sql.Stmt
}

func NewUserStorage(db *sql.DB) *UserStorage {
	return &UserStorage{
		DB:       db,
		Validate: validator.New(),
	}
}

func (u *UserStorage) Prepare(ctx context.Context) error {
	return u.prepareSqlStatements(ctx)
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
