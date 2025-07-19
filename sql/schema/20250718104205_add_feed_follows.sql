-- +goose Up
-- +goose StatementBegin
CREATE TABLE feed_follows (
    id uuid PRIMARY KEY,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    user_ID uuid REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    feed_ID uuid REFERENCES feeds (id) ON DELETE CASCADE NOT NULL,
    UNIQUE (user_ID, feed_ID)
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE feed_follows;

-- +goose StatementEnd
