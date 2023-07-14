-- name: CreateUser :one
Insert Into users(id, created_at, updated_at, name, api_key)
values ($1, $2, $3, $4, encode(sha256(random()::text::bytea), 'hex')) RETURNING *;

-- name: GetUserByApiKey :one
SELECT * FROM users WHERE api_key = $1;