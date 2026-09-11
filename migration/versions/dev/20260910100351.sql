-- Modify "users" table
ALTER TABLE "public"."users" DROP COLUMN "middle_name", DROP COLUMN "profile_image_file_id";
-- Drop "auth_sessions" table
DROP TABLE "public"."auth_sessions";
-- Drop "user_devices" table
DROP TABLE "public"."user_devices";
