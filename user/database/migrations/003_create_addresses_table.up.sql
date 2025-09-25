CREATE TABLE addresses(
    id BIGSERIAL PRIMARY KEY,
    auth_id BIGINT NOT NULL,
    receiver VARCHAR NOT NULL,
    phone VARCHAR NOT NULL,
    label VARCHAR NOT NULL,
    notes TEXT,
    is_primary BOOLEAN NOT NULL,
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
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index to optimize queries for addresses aggregated by auth_id
CREATE INDEX idx_addresses_auth_id ON addresses(auth_id);

-- Index to optimize queries for addresses by location point
CREATE INDEX idx_addresses_location ON addresses USING GIST (location);