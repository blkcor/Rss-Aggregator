-- name: CreateUser :one
Insert Into users(id,created_at,updated_at,name)
values ($1,$2,$3,$4)
RETURNING *;