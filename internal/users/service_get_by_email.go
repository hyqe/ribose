package users

import (
	"context"
	_ "embed"

	"github.com/hyqe/ribose/internal/fit/status"
	"github.com/lib/pq"
)

type GetByEmailRequest struct {
	Email string `json:"email"`
}
type GetByEmailResponse = User

func (u *UserStorage) GetByEmail(ctx context.Context, in *GetByEmailRequest) (*GetByEmailResponse, status.Status) {
	out := new(GetByEmailResponse)
	err := u.stmt_select_user_by_email.QueryRow(in.Email).Scan(out.Fields()...)
	if e, ok := err.(*pq.Error); ok {
		return nil, status.Pg(e)
	}
	return out, status.OK
}

//go:embed sql/select_user_by_email.sql
var sql_select_user_by_email string

func (u *UserStorage) prepareGetByEmail(ctx context.Context) error {
	stmt, err := u.DB.PrepareContext(ctx, sql_select_user_by_email)
	if err != nil {
		return err
	}
	u.stmt_select_user_by_email = stmt
	return nil
}
