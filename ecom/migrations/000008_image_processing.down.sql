ALTER TABLE image_jobs DROP CONSTRAINT IF EXISTS fk_image_jobs_image;
ALTER TABLE image_jobs DROP CONSTRAINT IF EXISTS image_jobs_status_check;
DROP TABLE IF EXISTS image_jobs;

ALTER TABLE images DROP CONSTRAINT IF EXISTS images_status_check;
ALTER TABLE images
    DROP COLUMN IF EXISTS object_key,
    DROP COLUMN IF EXISTS status,
    DROP COLUMN IF EXISTS mime,
    DROP COLUMN IF EXISTS width,
    DROP COLUMN IF EXISTS height,
    DROP COLUMN IF EXISTS processed_key,
    DROP COLUMN IF EXISTS error,
    ALTER COLUMN url SET NOT NULL;