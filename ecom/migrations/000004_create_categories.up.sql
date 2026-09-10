CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    name TEXT NOT NULL UNIQUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE listings
ADD COLUMN category_id UUID NOT NULL;

ALTER TABLE listings
ADD CONSTRAINT fk_listing_category
FOREIGN KEY (category_id)
REFERENCES categories(id);