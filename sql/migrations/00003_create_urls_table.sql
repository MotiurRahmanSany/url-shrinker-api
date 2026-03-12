-- +goose Up
-- +goose StatementBegin

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE urls (
    id           BIGSERIAL PRIMARY KEY,
    short_code   VARCHAR(10) UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    user_id      UUID REFERENCES users(id) ON DELETE SET NULL,
    is_active    BOOLEAN NOT NULL DEFAULT true,
    expires_at   TIMESTAMPTZ,          -- NULL = never expires
    max_clicks   INTEGER,              -- NULL = unlimited
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ON urls(short_code);
CREATE INDEX ON urls(user_id);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS urls;

-- +goose StatementEnd