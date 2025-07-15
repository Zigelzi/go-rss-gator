-- +goose Up
-- +goose StatementBegin
ALTER TABLE feeds
ALTER COLUMN user_ID SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE feeds
ALTER COLUMN user_ID DROP NOT NULL; 
-- +goose StatementEnd
