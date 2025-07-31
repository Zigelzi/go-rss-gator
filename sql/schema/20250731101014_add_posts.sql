-- +goose Up
-- +goose StatementBegin
CREATE TABLE posts (
    id uuid PRIMARY KEY,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    title text NOT NULL,
    description text,
    url text NOT NULL UNIQUE,
    published_at timestamp NOT NULL,
    feed_ID uuid NOT NULL REFERENCES feeds ON DELETE CASCADE 
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE posts;
-- +goose StatementEnd
