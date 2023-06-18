package users

import (
	"context"
	_ "embed"

	"github.com/google/uuid"
	"github.com/hyqe/ribose/internal/fit/status"
	"github.com/lib/pq"
)

type GetByUUIDRequest struct {
	UUID uuid.UUID `json:"uuid"`
}
type GetByUUIDResponse = User

func (s *Service) GetByUUID(ctx context.Context, in GetByUUIDRequest) (GetByUUIDResponse, status.Status) {
	var out GetByUUIDResponse
	err := s.stmt_select_user_by_uuid.QueryRow(in.UUID).Scan(out.Fields()...)
	if e, ok := err.(*pq.Error); ok {
		return out, status.Pg(e)
	}
	return out, status.OK
}

//go:embed sql/select_user_by_uuid.sql
var sql_select_user_by_uuid string

func (s *Service) prepareGetByUUID(ctx context.Context) error {
	stmt, err := s.DB.PrepareContext(ctx, sql_select_user_by_uuid)
	if err != nil {
		return err
	}
	s.stmt_select_user_by_uuid = stmt
	return nil
}
