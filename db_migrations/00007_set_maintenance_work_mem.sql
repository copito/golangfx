-- +goose NO TRANSACTION
-- +goose Up
ALTER SYSTEM SET maintenance_work_mem = '128MB';

-- +goose Down
ALTER SYSTEM RESET maintenance_work_mem;