-- +migrate Up
CREATE TABLE profiles (
    account_id  UUID PRIMARY KEY,
    username    VARCHAR(32) NOT NULL UNIQUE,
    official    BOOLEAN NOT NULL DEFAULT FALSE,
    pseudonym   VARCHAR(128),
    description VARCHAR(255),
    avatar      TEXT,

    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +migrate Down
DROP TABLE IF EXISTS profiles CASCADE;
DROP TABLE IF EXISTS outbox_events CASCADE;
DROP TABLE IF EXISTS inbox_events CASCADE;

DROP TYPE IF EXISTS outbox_event_status;
DROP TYPE IF EXISTS inbox_event_status;
