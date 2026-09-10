ALTER TABLE listings
DROP CONSTRAINT IF EXISTS fk_listing_category;

ALTER TABLE listings
DROP COLUMN IF EXISTS category_id;

DROP TABLE IF EXISTS categories;