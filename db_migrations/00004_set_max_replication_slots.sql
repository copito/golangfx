-- +goose NO TRANSACTION
-- +goose Up
ALTER SYSTEM SET max_replication_slots = 10;

-- +goose Down
ALTER SYSTEM RESET max_replication_slots;