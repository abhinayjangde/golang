ALTER TABLE listings
    DROP CONSTRAINT IF EXISTS fk_listings_user;

ALTER TABLE listings
    DROP COLUMN IF EXISTS user_id,
    DROP COLUMN IF EXISTS updated_at;