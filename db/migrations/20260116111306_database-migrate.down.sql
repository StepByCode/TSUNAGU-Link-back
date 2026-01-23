-- reverse: drop "users" table
CREATE TABLE "public"."users" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "email" character varying(255) NOT NULL,
  "name" character varying(255) NOT NULL,
  "password" character varying(255) NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamp NULL,
  PRIMARY KEY ("id")
);
CREATE INDEX "idx_users_deleted_at" ON "public"."users" ("deleted_at");
CREATE INDEX "idx_users_email" ON "public"."users" ("email") WHERE (deleted_at IS NULL);
CREATE UNIQUE INDEX "users_email_key" ON "public"."users" ("email");
-- reverse: create "profiles" table
DROP TABLE "profiles"."profiles";
-- reverse: create index "users_email_key" to table: "users"
DROP INDEX "users"."users_email_key";
-- reverse: create "users" table
DROP TABLE "users"."users";
-- reverse: create "locations" table
DROP TABLE "master_data"."locations";
-- reverse: create "grades" table
DROP TABLE "master_data"."grades";
-- reverse: Add new schema named "users"
DROP SCHEMA "users" CASCADE;
-- reverse: Add new schema named "profiles"
DROP SCHEMA "profiles" CASCADE;
-- reverse: Add new schema named "master_data"
DROP SCHEMA "master_data" CASCADE;
-- reverse: Add new schema named "attendance_logs"
DROP SCHEMA "attendance_logs" CASCADE;
