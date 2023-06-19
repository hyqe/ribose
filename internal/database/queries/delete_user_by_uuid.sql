-- name: DeleteUser :exec
DELETE 
FROM users 
WHERE 
    uuid = $1;