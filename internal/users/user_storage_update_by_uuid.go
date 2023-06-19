package users

import (
	"context"
	_ "embed"
	"log"

	"github.com/google/uuid"
	"github.com/hyqe/ribose/internal/fit/codes"
	"github.com/hyqe/ribose/internal/fit/status"
	"github.com/lib/pq"
)

type UpdateByUUIDRequest = struct {
	UUID uuid.UUID `json:"uuid"`
	User `json:"user"`
}
type UpdateByUUIDResponse = User

func (u *UserStorage) UpdateByUUID(ctx context.Context, in *UpdateByUUIDRequest) (*UpdateByUUIDResponse, status.Status) {
	out := new(UpdateByUUIDResponse)
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

	err = tx.QueryRowContext(ctx, sql_update_user_by_uuid, in.UUID, in.User.Email).Scan(out.Fields()...)
	if e, ok := err.(*pq.Error); ok {
		return nil, status.Pg(e)
	}

	err = tx.Commit()
	if err != nil {
		return nil, status.New(codes.Internal, err)
	}

	return out, status.OK

}

//go:embed sql/update_user_by_uuid.sql
var sql_update_user_by_uuid string
