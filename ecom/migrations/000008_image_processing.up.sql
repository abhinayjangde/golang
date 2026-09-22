ALTER TABLE images
    ALTER COLUMN url DROP NOT NULL,
    ADD COLUMN object_key TEXT NOT NULL DEFAULT '',
    ADD COLUMN status TEXT NOT NULL DEFAULT 'pending',
    ADD COLUMN mime TEXT,
    ADD COLUMN width INT,
    ADD COLUMN height INT,
    ADD COLUMN processed_key TEXT,
    ADD COLUMN error TEXT;

ALTER TABLE images
    ADD CONSTRAINT images_status_check
    CHECK (status IN ('pending', 'uploaded', 'processing', 'completed', 'failed'));

CREATE TABLE IF NOT EXISTS image_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    image_id UUID UNIQUE NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    attempts INT NOT NULL DEFAULT 0,
    last_error TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_image_jobs_status ON image_jobs(status);

ALTER TABLE image_jobs
    ADD CONSTRAINT image_jobs_status_check
    CHECK (status IN ('pending', 'processing', 'done', 'failed'));

ALTER TABLE image_jobs
    ADD CONSTRAINT fk_image_jobs_image
    FOREIGN KEY (image_id)
    REFERENCES images(id)
    ON DELETE CASCADE;