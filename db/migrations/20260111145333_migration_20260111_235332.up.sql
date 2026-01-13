-- modify "users" table
ALTER TABLE "public"."users" DROP CONSTRAINT "users_email_key", ALTER COLUMN "id" SET DEFAULT gen_random_uuid();
-- create index "users_email_key" to table: "users"
CREATE UNIQUE INDEX "users_email_key" ON "public"."users" ("email");
