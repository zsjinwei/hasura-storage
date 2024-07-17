ALTER TABLE "storage"."buckets" ADD COLUMN IF NOT EXISTS "upload_expiration" INT NOT NULL DEFAULT 600;

ALTER TABLE storage.buckets
    ADD CONSTRAINT upload_expiration_valid_range
        CHECK (upload_expiration >= 1 AND upload_expiration <= 604800);
