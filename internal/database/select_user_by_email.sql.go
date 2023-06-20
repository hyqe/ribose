// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: select_user_by_email.sql

package database

import (
	"context"
)

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, created_at, uuid, email
FROM users
WHERE
    email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Uuid,
		&i.Email,
	)
	return i, err
}