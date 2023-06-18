SELECT
    uuid,
    email
FROM users
WHERE
    email = $1;