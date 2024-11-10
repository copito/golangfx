-- +goose NO TRANSACTION
-- +goose Up
ALTER SYSTEM SET work_mem = '64MB';

-- +goose Down
ALTER SYSTEM RESET work_mem;