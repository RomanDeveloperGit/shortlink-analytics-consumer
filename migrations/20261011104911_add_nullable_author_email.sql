-- +goose Up
ALTER TABLE link_created_event
    MODIFY COLUMN author_email Nullable(String);

ALTER TABLE link_created_event
    UPDATE author_email = NULL
    WHERE author_email = ''
    SETTINGS mutations_sync = 1;

-- +goose Down
