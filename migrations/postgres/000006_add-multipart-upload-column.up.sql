ALTER TABLE "storage"."files" ADD COLUMN IF NOT EXISTS "object_key" TEXT;
ALTER TABLE "storage"."files" ADD COLUMN IF NOT EXISTS "chunk_size" INT;
ALTER TABLE "storage"."files" ADD COLUMN IF NOT EXISTS "chunk_count" INT;
ALTER TABLE "storage"."files" ADD COLUMN IF NOT EXISTS "upload_id" TEXT;
