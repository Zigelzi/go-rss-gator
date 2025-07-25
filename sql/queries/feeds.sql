-- name: CreateFeed :one
INSERT INTO
    feeds (id, created_at, updated_at, name, url, user_id)
VALUES
    ($1, $2, $3, $4, $5, $6)
RETURNING
    *;

-- name: GetFeeds :many
SELECT
    f.name AS feed_name,
    f.url AS feed_url,
    u.name AS user_name
FROM
    feeds f
    LEFT JOIN users u ON u.id = f.user_id;

-- name: GetFeedByURL :one
SELECT
    *
from
    feeds
WHERE
    feeds.url = $1;

-- name: FollowFeed :one
WITH
    inserted_feed_follow AS (
        INSERT INTO
            feed_follows AS ff (id, created_at, updated_at, user_ID, feed_ID)
        VALUES
            ($1, $2, $3, $4, $5)
        RETURNING
            *
    )
SELECT
    iff.*,
    f.name as feed_name,
    u.name as user_name
FROM
    inserted_feed_follow iff
    INNER JOIN feeds f ON f.id = iff.feed_ID
    INNER JOIN users u ON u.id = iff.user_ID;

-- name: GetUserFeedFollows :many
SELECT
    ff.*,
    f.name as feed_name,
    u.name as user_name
FROM
    feed_follows ff
    INNER JOIN feeds f ON f.id = ff.feed_ID
    INNER JOIN users u ON u.id = ff.user_ID
WHERE
    ff.user_ID = $1;

-- name: UnfollowFeed :exec
DELETE FROM feed_follows
WHERE user_ID=$1 AND feed_ID=$2;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET updated_at=$2,
    last_fetched_at=$3
WHERE id=$1;