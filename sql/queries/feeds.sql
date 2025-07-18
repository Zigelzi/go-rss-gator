-- name: CreateFeed :one
INSERT INTO feeds (
    id,
    created_at,
    updated_at,
    name,
    url,
    user_id
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
) RETURNING *;

-- name: GetFeeds :many
SELECT f.name as feed_name, f.url as feed_url, u.name as user_name FROM feeds f
LEFT JOIN users u ON u.id = f.user_id;