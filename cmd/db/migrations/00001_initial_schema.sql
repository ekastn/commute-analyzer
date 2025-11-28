-- +goose Up
CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE users (
    user_id    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id  TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE commutes (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id        UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    name           TEXT,
    home_point     GEOMETRY(POINT, 4326) NOT NULL,
    office_point   GEOMETRY(POINT, 4326) NOT NULL,
    route_geometry GEOMETRY(LINESTRING, 4326) NOT NULL,
    distance_km    DOUBLE PRECISION NOT NULL,
    duration_min   DOUBLE PRECISION NOT NULL,
    vehicle        TEXT NOT NULL CHECK (vehicle IN ('car', 'motorcycle')),
    fuel_price     INTEGER NOT NULL,
    days_per_week  SMALLINT NOT NULL CHECK (days_per_week BETWEEN 1 AND 7),
    annual_cost    BIGINT NOT NULL,
    annual_minutes BIGINT NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_commutes_user_id ON commutes(user_id);
CREATE INDEX idx_commutes_route ON commutes USING GIST(route_geometry);

-- +goose Down
DROP TABLE commutes;
DROP TABLE users;
