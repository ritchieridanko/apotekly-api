CREATE TABLE addresses(
    id BIGSERIAL PRIMARY KEY,
    auth_id BIGINT NOT NULL,
    is_primary BOOLEAN DEFAULT FALSE,
    street TEXT,
    subdistrict TEXT,
    district TEXT,
    city TEXT,
    province TEXT,
    postal_code TEXT,
    country TEXT,
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    location GEOGRAPHY(Point, 4326),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- Index to optimize queries for addresses by location point
CREATE INDEX idx_addresses_location ON addresses USING GIST (location);