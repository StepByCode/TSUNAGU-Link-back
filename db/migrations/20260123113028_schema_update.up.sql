-- create enum type "attendance_type"
CREATE TYPE "attendance_logs"."attendance_type" AS ENUM ('check-in', 'check-out');
-- modify "users" table
ALTER TABLE "users"."users" ADD COLUMN "zitadel_id" uuid NULL, ADD COLUMN "created_at" timestamptz NOT NULL DEFAULT now(), ADD COLUMN "updated_at" timestamptz NOT NULL DEFAULT now();
-- create index "users_zitadel_id_key" to table: "users"
CREATE UNIQUE INDEX "users_zitadel_id_key" ON "users"."users" ("zitadel_id");
-- create "attendance_logs" table
CREATE TABLE "attendance_logs"."attendance_logs" (
  "id" bigserial NOT NULL,
  "user_id" uuid NOT NULL,
  "location_id" integer NOT NULL,
  "type" "attendance_logs"."attendance_type" NOT NULL,
  "timestamp" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "attendance_logs_location_id_fkey" FOREIGN KEY ("location_id") REFERENCES "master_data"."locations" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "attendance_logs_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "attendance_logs_timestamp_idx" to table: "attendance_logs"
CREATE INDEX "attendance_logs_timestamp_idx" ON "attendance_logs"."attendance_logs" ("timestamp");
-- create index "attendance_logs_user_id_idx" to table: "attendance_logs"
CREATE INDEX "attendance_logs_user_id_idx" ON "attendance_logs"."attendance_logs" ("user_id");
-- create "organizations" table
CREATE TABLE "master_data"."organizations" (
  "id" serial NOT NULL,
  "name" text NOT NULL,
  PRIMARY KEY ("id")
);
-- modify "profiles" table
ALTER TABLE "profiles"."profiles" ADD COLUMN "display_id" text NOT NULL, ADD COLUMN "nfc_serial" text NULL, ADD COLUMN "org_id" integer NOT NULL, ADD CONSTRAINT "profiles_org_id_fkey" FOREIGN KEY ("org_id") REFERENCES "master_data"."organizations" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- create index "profiles_display_id_key" to table: "profiles"
CREATE UNIQUE INDEX "profiles_display_id_key" ON "profiles"."profiles" ("display_id");
-- create index "profiles_nfc_serial_key" to table: "profiles"
CREATE UNIQUE INDEX "profiles_nfc_serial_key" ON "profiles"."profiles" ("nfc_serial");
