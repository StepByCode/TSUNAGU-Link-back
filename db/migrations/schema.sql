-- 1. UUID機能の準備
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 2. master_data スキーマとテーブル
CREATE SCHEMA master_data;
CREATE TABLE master_data.grades (id SERIAL PRIMARY KEY, label TEXT NOT NULL);
CREATE TABLE master_data.organizations (id SERIAL PRIMARY KEY, name TEXT NOT NULL);
CREATE TABLE master_data.locations (id SERIAL PRIMARY KEY, name TEXT NOT NULL);

-- 3. users スキーマとテーブル
CREATE SCHEMA users;
CREATE TABLE users.users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    zitadel_id UUID UNIQUE,
    email TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 4. profiles スキーマとテーブル
CREATE SCHEMA profiles;
CREATE TABLE profiles.profiles (
    user_id UUID PRIMARY KEY REFERENCES users.users(id) ON DELETE CASCADE,
    display_id TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    nfc_serial TEXT UNIQUE,
    grade_id INTEGER REFERENCES master_data.grades(id),
    org_id INTEGER REFERENCES master_data.organizations(id)
);

-- 5. attendance_logs スキーマとテーブル
CREATE SCHEMA attendance_logs;
CREATE TABLE attendance_logs.attendance_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users.users(id),
    location_id INTEGER NOT NULL REFERENCES master_data.locations(id),
    type TEXT NOT NULL CHECK (type IN ('check-in', 'check-out')),
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);