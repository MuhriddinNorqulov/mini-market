-- Modify "users" table
ALTER TABLE "public"."users" DROP COLUMN "phone_number", ADD COLUMN "username" character varying(32) NULL;
-- Create index "idx_users_username" to table: "users"
CREATE UNIQUE INDEX "idx_users_username" ON "public"."users" ("username");
