-- +goose Up
ALTER TABLE link_created_event ADD COLUMN author_email String;

-- +goose Down
ALTER TABLE link_created_event DROP COLUMN author_email;
