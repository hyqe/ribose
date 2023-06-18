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

func (s *Service) DeleteByUUID(ctx context.Context, in DeleteByUUIDRequest) (DeleteByUUIDResponse, status.Status) {
	var out DeleteByUUIDResponse
	err := s.stmt_select_user_by_email.QueryRow(in.UUID).Err()
	if e, ok := err.(*pq.Error); ok {
		return out, status.Pg(e)
	}
	return out, status.OK
}

//go:embed sql/delete_user_by_uuid.sql
var sql_delete_user_by_uuid string

func (s *Service) prepareDeleteByUUID(ctx context.Context) error {
	stmt, err := s.DB.PrepareContext(ctx, sql_delete_user_by_uuid)
	if err != nil {
		return err
	}
	s.stmt_delete_user_by_uuid = stmt
	return nil
}
