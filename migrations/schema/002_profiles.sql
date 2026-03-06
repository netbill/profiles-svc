-- +migrate Up
CREATE TABLE accounts (
    id         UUID        PRIMARY KEY,
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
    username    VARCHAR(32) NOT NULL UNIQUE,
    pseudonym   VARCHAR(128),
    description VARCHAR(255),
    avatar_key  TEXT,
    version     INTEGER NOT NULL DEFAULT 1 CHECK ( version > 0 ),

    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE tombstones (
    id           UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    entity_type  VARCHAR(64) NOT NULL,
    entity_id    UUID        NOT NULL,
    deleted_at   TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE (entity_type, entity_id)
);

-- +migrate Down
DROP TABLE IF EXISTS accounts CASCADE;
DROP TABLE IF EXISTS profiles CASCADE;
DROP TABLE IF EXISTS tombstones CASCADE;
