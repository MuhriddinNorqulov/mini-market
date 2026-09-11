-- Modify "orders" table
ALTER TABLE "public"."orders" ADD CONSTRAINT "chk_orders_cancel_reason" CHECK ((cancel_reason)::text = ANY ((ARRAY['user_requested'::character varying, 'timeout'::character varying])::text[])), ADD CONSTRAINT "chk_orders_status" CHECK ((status)::text = ANY ((ARRAY['pending'::character varying, 'confirmed'::character varying, 'cancelled'::character varying])::text[]));
