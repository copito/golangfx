-- +goose NO TRANSACTION
-- +goose Up
ALTER SYSTEM SET wal_level = 'logical';

-- +goose Down
ALTER SYSTEM RESET wal_level;