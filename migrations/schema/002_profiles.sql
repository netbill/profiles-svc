-- +migrate Up
CREATE TABLE accounts (
    id         UUID        PRIMARY KEY,
    username   VARCHAR(32) NOT NULL UNIQUE,
    role       VARCHAR     NOT NULL,
    version    INTEGER     NOT NULL DEFAULT 1 CHECK ( version > 0 ),

    source_created_at  TIMESTAMPTZ NOT NULL,
    source_updated_at  TIMESTAMPTZ NOT NULL,
    replica_created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    replica_updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at         TIMESTAMPTZ
);

CREATE TABLE profiles (
    account_id  UUID PRIMARY KEY REFERENCES accounts(id) ON DELETE CASCADE,
    username    VARCHAR(32) NOT NULL UNIQUE,
    pseudonym   VARCHAR(128),
    description VARCHAR(255),
    avatar_key  TEXT,
    version     INTEGER NOT NULL DEFAULT 1 CHECK ( version > 0 ),

    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMPTZ
);

-- +migrate Down
DROP TABLE IF EXISTS accounts CASCADE;
DROP TABLE IF EXISTS profiles CASCADE;
