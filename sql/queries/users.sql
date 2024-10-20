-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email,hashed_password,is_chirpy_red)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    false
)
RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: UpdateUser :one
UPDATE users
SET hashed_password = $1, email = $2
WHERE id = $3
RETURNING *;

-- name: UpgradeRedUser :one
UPDATE users
SET is_chirpy_red = $1
WHERE id = $2
RETURNING *;