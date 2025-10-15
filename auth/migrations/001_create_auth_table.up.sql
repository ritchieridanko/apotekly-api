CREATE TABLE auth(
    auth_id BIGSERIAL PRIMARY KEY,

    -- Primary
    email VARCHAR NOT NULL,
    password VARCHAR,
    role SMALLINT NOT NULL,

    -- Secondary
    is_verified BOOLEAN DEFAULT FALSE,
    email_changed_at TIMESTAMPTZ,
    password_changed_at TIMESTAMPTZ,

    -- Metadata
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- Enforce uniqueness of email for active (not soft-deleted) records
CREATE UNIQUE INDEX idx_auth_unique_email ON auth(email) WHERE deleted_at IS NULL;

-- Index to optimize queries for verified and active records by email
CREATE INDEX idx_auth_active ON auth(email) WHERE is_verified = TRUE AND deleted_at IS NULL;