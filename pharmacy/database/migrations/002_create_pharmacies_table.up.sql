CREATE TABLE pharmacies(
    auth_id BIGINT PRIMARY KEY,
    pharmacy_public_id UUID UNIQUE NOT NULL, -- public-facing id
    pharmacy_id BIGSERIAL,

    -- Identity / Ownership
    name VARCHAR NOT NULL,
    legal_name VARCHAR, -- registered business/legal name
    description TEXT,
    license_number VARCHAR NOT NULL,
    license_authority VARCHAR NOT NULL,
    license_expiry DATE,

    -- Contact
    email VARCHAR,
    phone VARCHAR,
    website VARCHAR,

    -- Location
    country TEXT NOT NULL,
    admin_level_1 TEXT,
    admin_level_2 TEXT,
    admin_level_3 TEXT,
    admin_level_4 TEXT,
    street TEXT NOT NULL,
    postal_code TEXT NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    location GEOGRAPHY(Point, 4326),

    -- Media
    logo VARCHAR,

    -- Business Information
    opening_hours JSONB NOT NULL, -- e.g., { "mon": ["08:00-20:00"], "sun": [] }

    -- Operational Status
    status VARCHAR DEFAULT 'ACTIVE', -- ACTIVE, SUSPENDED, CLOSE
    is_active BOOLEAN DEFAULT TRUE,

    -- Metadata
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- Enforce uniqueness of license_number and license_authority for active (not soft-deleted) records
CREATE UNIQUE INDEX idx_pharmacies_unique_license
ON pharmacies(license_number, license_authority)
WHERE deleted_at IS NULL;

-- Index to optimize queries for active pharmacies by auth_id
CREATE INDEX idx_pharmacies_active
ON pharmacies(auth_id)
WHERE deleted_at IS NULL;

-- Index to optimize queries for active pharmacies by location point
CREATE INDEX idx_pharmacies_location
ON pharmacies USING GIST (location)
WHERE is_active = TRUE AND deleted_at IS NULL;

-- Index to optimize queries for text search
CREATE INDEX idx_pharmacies_text_search
ON pharmacies USING GIN (to_tsvector('simple', name || ' ' || coalesce(description, '')))
WHERE is_active = TRUE AND deleted_at IS NULL;