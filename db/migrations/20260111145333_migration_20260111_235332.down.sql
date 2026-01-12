-- reverse: create index "users_email_key" to table: "users"
DROP INDEX "public"."users_email_key";
-- reverse: modify "users" table
ALTER TABLE "public"."users" ALTER COLUMN "id" SET DEFAULT public.uuid_generate_v4(), ADD CONSTRAINT "users_email_key" UNIQUE ("email");
