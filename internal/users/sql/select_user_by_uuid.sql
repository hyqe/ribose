SELECT
    uuid,
    email
FROM users
WHERE
    uuid = $1;