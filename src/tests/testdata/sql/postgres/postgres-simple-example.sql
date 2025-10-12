-- Database initialization script
-- This script creates tables and inserts fake data for dev testing

-- Create schemas if they dont exist
CREATE SCHEMA IF NOT EXISTS example_schema;

-- Create simple table
CREATE TABLE IF NOT EXISTS example_schema.example_table (
    id  UUID NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'UTC')
);

-- Insert fake data into the table
INSERT INTO example_schema.example_table(id, name) VALUES
    ('577bd4e4-30ed-4129-beb8-1168cc1b45f1', 'Example Name 1'),
    ('ceca9623-a94e-412b-a08c-57931c26ee0d', 'Example Name 2'),
    ('d7e7d0ba-6998-4ab4-b08b-e811a855c6ae', 'Example Name 3')