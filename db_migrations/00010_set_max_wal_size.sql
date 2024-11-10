-- +goose NO TRANSACTION
-- +goose Up
ALTER SYSTEM SET max_wal_size = '1GB';

-- +goose Down
ALTER SYSTEM RESET max_wal_size;