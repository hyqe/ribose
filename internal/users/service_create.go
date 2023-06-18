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
	Email string `json:"email"`
}
type CreateResponse = User

func (s *Service) Create(ctx context.Context, in CreateRequest) (CreateResponse, status.Status) {
	var out CreateResponse
	tx, err := s.DB.Begin()
	if err != nil {
		return out, status.Newf(codes.Internal, "failed to begin transaction: %v", err)
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
		return out, status.Pg(e)
	}

	err = tx.Commit()
	if err != nil {
		return out, status.New(codes.Internal, err)
	}

	return out, status.OK
}

//go:embed sql/insert_user.sql
var sql_insert_user string
