-- +goose NO TRANSACTION
-- +goose Up
ALTER SYSTEM SET max_wal_senders = 10;

-- +goose Down
ALTER SYSTEM RESET max_wal_senders;