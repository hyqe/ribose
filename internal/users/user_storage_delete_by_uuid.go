package users

import (
	"context"
	_ "embed"

	"github.com/google/uuid"
	"github.com/hyqe/ribose/internal/fit/status"
	"github.com/lib/pq"
)

type DeleteByUUIDRequest struct {
	UUID uuid.UUID `json:"uuid"`
}
type DeleteByUUIDResponse struct{}

func (u *UserStorage) DeleteByUUID(ctx context.Context, in *DeleteByUUIDRequest) (*DeleteByUUIDResponse, status.Status) {
	out := new(DeleteByUUIDResponse)
	err := u.stmt_delete_user_by_uuid.QueryRow(in.UUID).Err()
	if e, ok := err.(*pq.Error); ok {
		return nil, status.Pg(e)
	}
	return out, status.OK
}

//go:embed sql/delete_user_by_uuid.sql
var sql_delete_user_by_uuid string

func (u *UserStorage) prepareDeleteByUUID(ctx context.Context) error {
	stmt, err := u.DB.PrepareContext(ctx, sql_delete_user_by_uuid)
	if err != nil {
		return err
	}
	u.stmt_delete_user_by_uuid = stmt
	return nil
}
