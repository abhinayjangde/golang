ALTER TABLE listings
    ADD COLUMN user_id UUID,
    ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

ALTER TABLE listings
    ADD CONSTRAINT fk_listings_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE;
