package users

import (
	"context"

	"github.com/google/uuid"
	"github.com/hyqe/ribose/internal/database"
	"github.com/hyqe/ribose/internal/fit/codes"
	"github.com/hyqe/ribose/internal/fit/status"
)

type Service struct {
	db *database.Queries
}

func NewService(db *database.Queries) *Service {
	return &Service{
		db: db,
	}
}

type CreateRequest struct {
	Email string `json:"email" validate:"email" example:"foo@example.com"`
}
type CreateResponse = User

func (s *Service) Create(ctx context.Context, in *CreateRequest) (*CreateResponse, status.Status) {
	switch u, err := s.db.CreateUsers(ctx, in.Email); err {
	case nil:
		return &CreateResponse{
			UUID:  u.Uuid,
			Email: u.Email,
		}, status.OK
	default:
		return nil, status.New(codes.Internal, err)
	}
}

type UpdateByUUIDRequest = struct {
	UUID uuid.UUID `json:"uuid"`
	User `json:"user"`
}
type UpdateByUUIDResponse = User

func (s *Service) UpdateByUUID(ctx context.Context, in *UpdateByUUIDRequest) (*UpdateByUUIDResponse, status.Status) {
	switch u, err := s.db.UpdateUserByUUID(ctx, database.UpdateUserByUUIDParams{
		Uuid:  in.UUID,
		Email: in.User.Email,
	}); err {
	case nil:
		return &UpdateByUUIDResponse{
			UUID:  u.Uuid,
			Email: u.Email,
		}, status.OK
	default:
		return nil, status.New(codes.Internal, err)
	}
}

type DeleteByUUIDRequest struct {
	UUID uuid.UUID `json:"uuid"`
}
type DeleteByUUIDResponse struct{}

func (s *Service) DeleteByUUID(ctx context.Context, in *DeleteByUUIDRequest) (*DeleteByUUIDResponse, status.Status) {
	switch err := s.db.DeleteUser(ctx, in.UUID); err {
	case nil:
		return &DeleteByUUIDResponse{}, status.OK
	default:
		return nil, status.New(codes.Internal, err)
	}
}

type GetByEmailRequest struct {
	Email string `json:"email" validate:"email" example:"foo@example.com"`
}
type GetByEmailResponse = User

func (s *Service) GetByEmail(ctx context.Context, in *GetByEmailRequest) (*GetByEmailResponse, status.Status) {
	switch u, err := s.db.GetUserByEmail(ctx, in.Email); err {
	case nil:
		return &GetByEmailResponse{
			UUID:  u.Uuid,
			Email: u.Email,
		}, status.OK
	default:
		return nil, status.New(codes.Internal, err)
	}
}

type GetByUUIDRequest struct {
	UUID uuid.UUID `json:"uuid"`
}
type GetByUUIDResponse = User

func (s *Service) GetByUUID(ctx context.Context, in *GetByUUIDRequest) (*GetByUUIDResponse, status.Status) {
	switch u, err := s.db.GetUserByUUID(ctx, in.UUID); err {
	case nil:
		return &GetByUUIDResponse{
			UUID:  u.Uuid,
			Email: u.Email,
		}, status.OK
	default:
		return nil, status.New(codes.Internal, err)
	}
}
