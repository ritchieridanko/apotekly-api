CREATE TABLE users(
    auth_id BIGINT PRIMARY KEY,
    user_id UUID UNIQUE NOT NULL,

    -- Personal Information
    name VARCHAR NOT NULL,
    bio TEXT,
    sex VARCHAR,
    birthdate DATE,
    phone VARCHAR,
    profile_picture VARCHAR,

    -- Metadata
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);