-- +goose NO TRANSACTION
-- +goose Up
ALTER SYSTEM SET checkpoint_timeout = '10min';

-- +goose Down
ALTER SYSTEM RESET checkpoint_timeout;