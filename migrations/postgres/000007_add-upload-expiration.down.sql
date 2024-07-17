ALTER TABLE storage.buckets
    DROP CONSTRAINT upload_expiration_valid_range;

ALTER TABLE "storage"."buckets" DROP COLUMN IF EXISTS "upload_expiration";