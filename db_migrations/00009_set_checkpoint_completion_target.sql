-- +goose NO TRANSACTION
-- +goose Up
ALTER SYSTEM SET checkpoint_completion_target = 0.9;

-- +goose Down
ALTER SYSTEM RESET checkpoint_completion_target;