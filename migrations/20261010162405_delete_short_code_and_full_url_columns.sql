-- +goose Up
ALTER TABLE link_visited_event DROP COLUMN short_code;
ALTER TABLE link_visited_event DROP COLUMN full_url;

-- +goose Down
ALTER TABLE link_visited_event ADD COLUMN short_code String;
ALTER TABLE link_visited_event ADD COLUMN full_url String;
