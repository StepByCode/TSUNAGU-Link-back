-- reverse: create index "profiles_nfc_serial_key" to table: "profiles"
DROP INDEX "profiles"."profiles_nfc_serial_key";
-- reverse: create index "profiles_display_id_key" to table: "profiles"
DROP INDEX "profiles"."profiles_display_id_key";
-- reverse: modify "profiles" table
ALTER TABLE "profiles"."profiles" DROP CONSTRAINT "profiles_org_id_fkey", DROP COLUMN "org_id", DROP COLUMN "nfc_serial", DROP COLUMN "display_id";
-- reverse: create "organizations" table
DROP TABLE "master_data"."organizations";
-- reverse: create index "attendance_logs_user_id_idx" to table: "attendance_logs"
DROP INDEX "attendance_logs"."attendance_logs_user_id_idx";
-- reverse: create index "attendance_logs_timestamp_idx" to table: "attendance_logs"
DROP INDEX "attendance_logs"."attendance_logs_timestamp_idx";
-- reverse: create "attendance_logs" table
DROP TABLE "attendance_logs"."attendance_logs";
-- reverse: create index "users_zitadel_id_key" to table: "users"
DROP INDEX "users"."users_zitadel_id_key";
-- reverse: modify "users" table
ALTER TABLE "users"."users" DROP COLUMN "updated_at", DROP COLUMN "created_at", DROP COLUMN "zitadel_id";
-- reverse: create enum type "attendance_type"
DROP TYPE "attendance_logs"."attendance_type";
