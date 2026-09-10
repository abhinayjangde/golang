CREATE INDEX idx_listings_user_id
ON listings(user_id);

CREATE INDEX idx_listings_category_id
ON listings(category_id);

CREATE INDEX idx_images_listing_id
ON images(listing_id);