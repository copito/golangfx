-- +goose NO TRANSACTION
-- +goose Up
ALTER SYSTEM SET shared_buffers = '256MB';

-- +goose Down
ALTER SYSTEM RESET shared_buffers;