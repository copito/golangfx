-- +goose Up
-- +goose StatementBegin
-- Create a sample table with UUID primary key
CREATE TABLE
    IF NOT EXISTS db_logs (
        event_uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
        type VARCHAR(50) NOT NULL,
        message VARCHAR(100) UNIQUE NOT NULL,
        payload JSONB,
        created_at TIMESTAMPTZ DEFAULT NOW (),
    );

-- Convert to hypertable
SELECT
    public.create_hypertable (
        'db_logs',
        by_range ('created_at'),
        if_not_exists = > TRUE
    );

-- Add retention policy for 365 days
SELECT
    add_retention_policy ('db_logs', INTERVAL '365 days');

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
SELECT
    remove_retention_policy ('db_logs');

DROP TABLE IF EXISTS db_logs;

-- +goose StatementEnd