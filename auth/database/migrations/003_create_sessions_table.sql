CREATE TABLE sessions(
    id BIGSERIAL PRIMARY KEY,
    auth_id BIGINT NOT NULL,
    parent_id BIGINT,
    token VARCHAR UNIQUE NOT NULL,
    user_agent TEXT NOT NULL,
    ip_address TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    FOREIGN KEY (auth_id) REFERENCES auth(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES sessions(id) ON DELETE CASCADE
);

-- Index to optimize queries for records by auth_id
CREATE INDEX idx_sessions_auth_id ON sessions(auth_id);

-- Index to optimize queries for records with parent_id not null
CREATE INDEX idx_sessions_parent_id ON sessions(parent_id) WHERE parent_id IS NOT NULL;

-- Index to optimize queries for active (not revoked) records by token
CREATE INDEX idx_sessions_active ON sessions(token) WHERE revoked_at IS NULL;