package users

import (
	"github.com/google/uuid"
)

type User struct {
	UUID  uuid.UUID `json:"uuid"`
	Email string    `json:"email"`
}

func (u *User) Fields() []any {
	return []any{
		&u.UUID,
		&u.Email,
	}
}
