package users

import "github.com/google/uuid"

type User struct {
	UUID  uuid.UUID `json:"uuid"`
	Email string    `json:"email" validate:"email" example:"foo@example.com"`
}
