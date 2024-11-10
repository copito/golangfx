-- +goose Up

-- Enable the UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Enable pgcrypto extension for cryptographic functions if needed
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- +goose Down

-- Drop extensions
DROP EXTENSION IF EXISTS pgcrypto;
DROP EXTENSION IF EXISTS "uuid-ossp";

