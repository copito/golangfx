-- +goose NO TRANSACTION
-- +goose Up
ALTER SYSTEM SET min_wal_size = '80MB';

-- +goose Down
ALTER SYSTEM RESET min_wal_size;