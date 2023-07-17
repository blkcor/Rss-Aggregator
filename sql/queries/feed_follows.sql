-- name: CreateFeedFollows :one
Insert Into feed_follows(id, created_at, updated_at, user_id, feed_id)
values ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetFeedFollows :many
SELECT *
FROM feed_follows
WHERE user_id = $1;

-- name: DeleteFeedFollows :exec
DELETE FROM feed_follows WHERE user_id = $1 AND id = $2;