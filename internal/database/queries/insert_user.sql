-- name: CreateUsers :one
INSERT INTO users (
	email
) 
VALUES (
	$1
) 
RETURNING *;