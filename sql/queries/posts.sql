-- name: CreatePost :one
INSERT INTO
    posts (
        id,
        created_at,
        updated_at,
        title,
        description,
        url,
        published_at,
        feed_ID
    )
VALUES
    ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING
    *;

-- name: GetPostsForUser :many
SELECT
    p.*
FROM
    posts p
    INNER JOIN feeds f ON f.id = p.feed_ID
    INNER JOIN feed_follows ff ON ff.feed_id = f.id
WHERE
    ff.user_id = $1
ORDER BY
    p.published_at DESC
LIMIT
    $2;