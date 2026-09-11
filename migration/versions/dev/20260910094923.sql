-- Create "users" table
CREATE TABLE "public"."users" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "google_id" character varying(128) NULL,
  "email" character varying(256) NULL,
  "picture" character varying(256) NULL,
  "email_verified" boolean NOT NULL DEFAULT false,
  "phone_number" character varying(32) NULL,
  "first_name" character varying(64) NULL,
  "last_name" character varying(64) NULL,
  "middle_name" character varying(64) NULL,
  "profile_image_file_id" bigint NULL,
  "role" text NOT NULL DEFAULT 'USER',
  "password" text NULL,
  "password_updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_users_deleted_at" to table: "users"
CREATE INDEX "idx_users_deleted_at" ON "public"."users" ("deleted_at");
-- Create index "idx_users_email" to table: "users"
CREATE UNIQUE INDEX "idx_users_email" ON "public"."users" ("email");
-- Create index "idx_users_google_id" to table: "users"
CREATE UNIQUE INDEX "idx_users_google_id" ON "public"."users" ("google_id");
-- Create index "idx_users_phone_number" to table: "users"
CREATE UNIQUE INDEX "idx_users_phone_number" ON "public"."users" ("phone_number");
-- Create index "idx_users_profile_image_file_id" to table: "users"
CREATE INDEX "idx_users_profile_image_file_id" ON "public"."users" ("profile_image_file_id");
-- Create "user_devices" table
CREATE TABLE "public"."user_devices" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "user_id" bigint NOT NULL,
  "device_id" character varying(64) NOT NULL,
  "name" character varying(128) NOT NULL,
  "browser_family" character varying(32) NOT NULL,
  "platform_family" character varying(32) NOT NULL,
  "last_seen_at" timestamptz NOT NULL,
  "last_ip" character varying(64) NULL,
  "trusted_at" timestamptz NULL,
  "trust_expires_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_user_devices_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create index "idx_user_devices_deleted_at" to table: "user_devices"
CREATE INDEX "idx_user_devices_deleted_at" ON "public"."user_devices" ("deleted_at");
-- Create index "idx_user_devices_user_device" to table: "user_devices"
CREATE UNIQUE INDEX "idx_user_devices_user_device" ON "public"."user_devices" ("user_id", "device_id");
-- Create "auth_sessions" table
CREATE TABLE "public"."auth_sessions" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "session_id" character varying(64) NOT NULL,
  "user_id" bigint NOT NULL,
  "user_device_id" bigint NOT NULL,
  "refresh_token_hash" character varying(64) NOT NULL,
  "issued_at" timestamptz NOT NULL,
  "last_seen_at" timestamptz NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "revoked_at" timestamptz NULL,
  "revoked_reason" character varying(32) NULL,
  "ip" character varying(64) NULL,
  "user_agent" character varying(512) NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_auth_sessions_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT "fk_auth_sessions_user_device" FOREIGN KEY ("user_device_id") REFERENCES "public"."user_devices" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create index "idx_auth_sessions_deleted_at" to table: "auth_sessions"
CREATE INDEX "idx_auth_sessions_deleted_at" ON "public"."auth_sessions" ("deleted_at");
-- Create index "idx_auth_sessions_refresh_token_hash" to table: "auth_sessions"
CREATE UNIQUE INDEX "idx_auth_sessions_refresh_token_hash" ON "public"."auth_sessions" ("refresh_token_hash");
-- Create index "idx_auth_sessions_session_id" to table: "auth_sessions"
CREATE UNIQUE INDEX "idx_auth_sessions_session_id" ON "public"."auth_sessions" ("session_id");
-- Create index "idx_auth_sessions_user_active" to table: "auth_sessions"
CREATE INDEX "idx_auth_sessions_user_active" ON "public"."auth_sessions" ("user_id", "revoked_at", "expires_at");
-- Create index "idx_auth_sessions_user_device_id" to table: "auth_sessions"
CREATE INDEX "idx_auth_sessions_user_device_id" ON "public"."auth_sessions" ("user_device_id");
-- Create "products" table
CREATE TABLE "public"."products" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "name" text NOT NULL,
  "price" bigint NOT NULL,
  "stock_quantity" bigint NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_products_deleted_at" to table: "products"
CREATE INDEX "idx_products_deleted_at" ON "public"."products" ("deleted_at");
-- Create "orders" table
CREATE TABLE "public"."orders" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "idempotency_key" character varying(64) NOT NULL,
  "user_id" bigint NOT NULL,
  "total_price" bigint NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "status" character varying(32) NOT NULL DEFAULT 'pending',
  "cancel_reason" character varying(32) NULL,
  "completed_at" timestamptz NULL,
  "cancelled_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_orders_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE SET NULL
);
-- Create index "idx_order_idempotency_key" to table: "orders"
CREATE UNIQUE INDEX "idx_order_idempotency_key" ON "public"."orders" ("idempotency_key", "user_id");
-- Create index "idx_orders_deleted_at" to table: "orders"
CREATE INDEX "idx_orders_deleted_at" ON "public"."orders" ("deleted_at");
-- Create "order_items" table
CREATE TABLE "public"."order_items" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "order_id" bigint NOT NULL,
  "product_id" bigint NOT NULL,
  "quantity" bigint NOT NULL,
  "unit_price" bigint NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_order_items_product" FOREIGN KEY ("product_id") REFERENCES "public"."products" ("id") ON UPDATE NO ACTION ON DELETE SET NULL,
  CONSTRAINT "fk_orders_items" FOREIGN KEY ("order_id") REFERENCES "public"."orders" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_order_items_deleted_at" to table: "order_items"
CREATE INDEX "idx_order_items_deleted_at" ON "public"."order_items" ("deleted_at");
-- Create index "idx_order_items_order_id_product_id" to table: "order_items"
CREATE UNIQUE INDEX "idx_order_items_order_id_product_id" ON "public"."order_items" ("order_id", "product_id");
