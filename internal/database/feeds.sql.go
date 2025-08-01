// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: feeds.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createFeed = `-- name: CreateFeed :one
INSERT INTO
    feeds (id, created_at, updated_at, name, url, user_id)
VALUES
    ($1, $2, $3, $4, $5, $6)
RETURNING
    id, created_at, updated_at, name, url, user_id, last_fetched_at
`

type CreateFeedParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Url       string
	UserID    uuid.UUID
}

func (q *Queries) CreateFeed(ctx context.Context, arg CreateFeedParams) (Feed, error) {
	row := q.db.QueryRowContext(ctx, createFeed,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Name,
		arg.Url,
		arg.UserID,
	)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Url,
		&i.UserID,
		&i.LastFetchedAt,
	)
	return i, err
}

const followFeed = `-- name: FollowFeed :one
WITH
    inserted_feed_follow AS (
        INSERT INTO
            feed_follows AS ff (id, created_at, updated_at, user_ID, feed_ID)
        VALUES
            ($1, $2, $3, $4, $5)
        RETURNING
            id, created_at, updated_at, user_id, feed_id
    )
SELECT
    iff.id, iff.created_at, iff.updated_at, iff.user_id, iff.feed_id,
    f.name as feed_name,
    u.name as user_name
FROM
    inserted_feed_follow iff
    INNER JOIN feeds f ON f.id = iff.feed_ID
    INNER JOIN users u ON u.id = iff.user_ID
`

type FollowFeedParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	FeedID    uuid.UUID
}

type FollowFeedRow struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	FeedID    uuid.UUID
	FeedName  string
	UserName  string
}

func (q *Queries) FollowFeed(ctx context.Context, arg FollowFeedParams) (FollowFeedRow, error) {
	row := q.db.QueryRowContext(ctx, followFeed,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.UserID,
		arg.FeedID,
	)
	var i FollowFeedRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.FeedID,
		&i.FeedName,
		&i.UserName,
	)
	return i, err
}

const getFeedByURL = `-- name: GetFeedByURL :one
SELECT
    id, created_at, updated_at, name, url, user_id, last_fetched_at
from
    feeds
WHERE
    feeds.url = $1
`

func (q *Queries) GetFeedByURL(ctx context.Context, url string) (Feed, error) {
	row := q.db.QueryRowContext(ctx, getFeedByURL, url)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Url,
		&i.UserID,
		&i.LastFetchedAt,
	)
	return i, err
}

const getFeeds = `-- name: GetFeeds :many
SELECT
    f.name AS feed_name,
    f.url AS feed_url,
    u.name AS user_name
FROM
    feeds f
    LEFT JOIN users u ON u.id = f.user_id
`

type GetFeedsRow struct {
	FeedName string
	FeedUrl  string
	UserName sql.NullString
}

func (q *Queries) GetFeeds(ctx context.Context) ([]GetFeedsRow, error) {
	rows, err := q.db.QueryContext(ctx, getFeeds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFeedsRow
	for rows.Next() {
		var i GetFeedsRow
		if err := rows.Scan(&i.FeedName, &i.FeedUrl, &i.UserName); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getNextFeedToFetch = `-- name: GetNextFeedToFetch :one
SELECT id, created_at, last_fetched_at, url, name
FROM feeds
ORDER BY last_fetched_at asc NULLS FIRST, created_at ASC
LIMIT 1
`

type GetNextFeedToFetchRow struct {
	ID            uuid.UUID
	CreatedAt     time.Time
	LastFetchedAt sql.NullTime
	Url           string
	Name          string
}

func (q *Queries) GetNextFeedToFetch(ctx context.Context) (GetNextFeedToFetchRow, error) {
	row := q.db.QueryRowContext(ctx, getNextFeedToFetch)
	var i GetNextFeedToFetchRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.LastFetchedAt,
		&i.Url,
		&i.Name,
	)
	return i, err
}

const getUserFeedFollows = `-- name: GetUserFeedFollows :many
SELECT
    ff.id, ff.created_at, ff.updated_at, ff.user_id, ff.feed_id,
    f.name as feed_name,
    u.name as user_name
FROM
    feed_follows ff
    INNER JOIN feeds f ON f.id = ff.feed_ID
    INNER JOIN users u ON u.id = ff.user_ID
WHERE
    ff.user_ID = $1
`

type GetUserFeedFollowsRow struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	FeedID    uuid.UUID
	FeedName  string
	UserName  string
}

func (q *Queries) GetUserFeedFollows(ctx context.Context, userID uuid.UUID) ([]GetUserFeedFollowsRow, error) {
	rows, err := q.db.QueryContext(ctx, getUserFeedFollows, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUserFeedFollowsRow
	for rows.Next() {
		var i GetUserFeedFollowsRow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
			&i.FeedID,
			&i.FeedName,
			&i.UserName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const markFeedFetched = `-- name: MarkFeedFetched :exec
UPDATE feeds
SET updated_at=NOW(),
    last_fetched_at=NOW()
WHERE id=$1
`

func (q *Queries) MarkFeedFetched(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, markFeedFetched, id)
	return err
}

const unfollowFeed = `-- name: UnfollowFeed :exec
DELETE FROM feed_follows
WHERE user_ID=$1 AND feed_ID=$2
`

type UnfollowFeedParams struct {
	UserID uuid.UUID
	FeedID uuid.UUID
}

func (q *Queries) UnfollowFeed(ctx context.Context, arg UnfollowFeedParams) error {
	_, err := q.db.ExecContext(ctx, unfollowFeed, arg.UserID, arg.FeedID)
	return err
}
