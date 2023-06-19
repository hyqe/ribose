package users

import (
	"context"
	_ "embed"
	"log"

	"github.com/hyqe/ribose/internal/fit/codes"
	"github.com/hyqe/ribose/internal/fit/status"
	"github.com/lib/pq"
)

type CreateRequest struct {
	Email string `json:"email" validate:"email"`
}
type CreateResponse = User

func (u *UserStorage) Create(ctx context.Context, in *CreateRequest) (*CreateResponse, status.Status) {
	out := new(CreateResponse)
	tx, err := u.DB.Begin()
	if err != nil {
		return nil, status.Newf(codes.Internal, "failed to begin transaction: %v", err)
	}
	defer func() {
		if err != nil {
			if err = tx.Rollback(); err != nil {
				log.Printf("failed to rollback transaction: %v", err)
			}
		}
	}()

	err = tx.QueryRowContext(ctx, sql_insert_user, in.Email).Scan(out.Fields()...)
	if e, ok := err.(*pq.Error); ok {
		return nil, status.Pg(e)
	}

	err = tx.Commit()
	if err != nil {
		return nil, status.New(codes.Internal, err)
	}

	return out, status.OK
}

//go:embed sql/insert_user.sql
var sql_insert_user string
