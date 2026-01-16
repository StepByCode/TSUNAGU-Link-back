-- Add new schema named "attendance_logs"
CREATE SCHEMA "attendance_logs";
-- Add new schema named "master_data"
CREATE SCHEMA "master_data";
-- Add new schema named "profiles"
CREATE SCHEMA "profiles";
-- Add new schema named "users"
CREATE SCHEMA "users";
-- create "grades" table
CREATE TABLE "master_data"."grades" (
  "id" serial NOT NULL,
  "label" text NOT NULL,
  PRIMARY KEY ("id")
);
-- create "locations" table
CREATE TABLE "master_data"."locations" (
  "id" serial NOT NULL,
  "name" text NOT NULL,
  PRIMARY KEY ("id")
);
-- create "users" table
CREATE TABLE "users"."users" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "email" text NOT NULL,
  PRIMARY KEY ("id")
);
-- create index "users_email_key" to table: "users"
CREATE UNIQUE INDEX "users_email_key" ON "users"."users" ("email");
-- create "profiles" table
CREATE TABLE "profiles"."profiles" (
  "user_id" uuid NOT NULL,
  "name" text NOT NULL,
  "grade_id" integer NOT NULL,
  PRIMARY KEY ("user_id"),
  CONSTRAINT "profiles_grade_id_fkey" FOREIGN KEY ("grade_id") REFERENCES "master_data"."grades" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "profiles_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- drop "users" table
DROP TABLE "public"."users";
