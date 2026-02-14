-- +migrate Up
CREATE TABLE accounts (
    id         UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    username   VARCHAR(32) NOT NULL UNIQUE,
    role       VARCHAR     NOT NULL,
    version    INTEGER     NOT NULL DEFAULT 1 CHECK ( version > 0 ),

    source_created_at  TIMESTAMPTZ NOT NULL,
    source_updated_at  TIMESTAMPTZ NOT NULL,
    replica_created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    replica_updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE profiles (
    account_id  UUID PRIMARY KEY REFERENCES accounts(id) ON DELETE CASCADE,
    username    VARCHAR(32) NOT NULL UNIQUE REFERENCES accounts(username) ON DELETE CASCADE,
    official    BOOLEAN NOT NULL DEFAULT FALSE,
    pseudonym   VARCHAR(128),
    description VARCHAR(255),
    avatar      TEXT,
    version     INTEGER NOT NULL DEFAULT 1 CHECK ( version > 0 ),

    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +migrate Down
DROP TABLE IF EXISTS accounts CASCADE;
DROP TABLE IF EXISTS profiles CASCADE;
DROP TABLE IF EXISTS outbox_events CASCADE;
DROP TABLE IF EXISTS inbox_events CASCADE;

DROP TYPE IF EXISTS outbox_event_status;
DROP TYPE IF EXISTS inbox_event_status;
