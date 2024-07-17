ALTER TABLE "storage"."files" DROP COLUMN IF EXISTS "upload_id";
ALTER TABLE "storage"."files" DROP COLUMN IF EXISTS "chunk_count";
ALTER TABLE "storage"."files" DROP COLUMN IF EXISTS "chunk_size";
ALTER TABLE "storage"."files" DROP COLUMN IF EXISTS "object_key";
