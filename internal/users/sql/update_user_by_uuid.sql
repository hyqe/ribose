UPDATE users
SET
    email=$2
WHERE
    uuid=$1
RETURNING 
	uuid, 
	email;