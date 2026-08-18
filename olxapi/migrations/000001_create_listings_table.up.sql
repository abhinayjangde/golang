CREATE TABLE IF NOT EXISTS listing (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    description TEXT,
    price BIGINT NOT NULL,
    city TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);