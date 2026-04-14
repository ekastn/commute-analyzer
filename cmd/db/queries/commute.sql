-- name: GetCommute :one
SELECT 
    id, user_id, name, 
    ST_AsBinary(home_point) AS home_point, 
    ST_AsBinary(office_point) AS office_point, 
    distance_km, duration_min, vehicle, fuel_price, days_per_week, 
    annual_cost, annual_minutes, created_at, updated_at,
    ST_AsBinary(route_geometry) AS route_geometry
FROM commutes WHERE id = $1;

-- name: CreateCommute :one
INSERT INTO commutes (
    user_id, name, home_point, office_point, route_geometry,
    distance_km, duration_min, vehicle, fuel_price, days_per_week,
    annual_cost, annual_minutes
) VALUES (
    $1, $2, 
    ST_SetSRID(ST_GeomFromWKB(@home_point::bytea), 4326), 
    ST_SetSRID(ST_GeomFromWKB(@office_point::bytea), 4326), 
    ST_SetSRID(ST_GeomFromWKB(@route_geometry::bytea), 4326), 
    $3, $4, $5, $6, $7, $8, $9
)
RETURNING id, annual_cost, annual_minutes, created_at;

-- name: ListCommutesByUser :many
SELECT 
    id, user_id, name, 
    ST_AsBinary(home_point) AS home_point, 
    ST_AsBinary(office_point) AS office_point, 
    distance_km, duration_min, vehicle, fuel_price, days_per_week, 
    annual_cost, annual_minutes, created_at, updated_at,
    ST_AsBinary(route_geometry) AS route_geometry
FROM commutes WHERE user_id = $1 ORDER BY created_at DESC;

-- name: UpdateCommute :one
UPDATE commutes SET
    name = COALESCE(sqlc.narg(name), name),
    vehicle = COALESCE(sqlc.narg(vehicle), vehicle),
    fuel_price = COALESCE(sqlc.narg(fuel_price), fuel_price),
    days_per_week = COALESCE(sqlc.narg(days_per_week), days_per_week),
    home_point = CASE
        WHEN sqlc.narg(home_point)::bytea IS NOT NULL
        THEN ST_SetSRID(ST_GeomFromWKB(sqlc.narg(home_point)::bytea), 4326)
        ELSE home_point
    END,
    office_point = CASE
        WHEN sqlc.narg(office_point)::bytea IS NOT NULL
        THEN ST_SetSRID(ST_GeomFromWKB(sqlc.narg(office_point)::bytea), 4326)
        ELSE office_point
    END,
    route_geometry = CASE
        WHEN sqlc.narg(route_geometry)::bytea IS NOT NULL
        THEN ST_SetSRID(ST_GeomFromWKB(sqlc.narg(route_geometry)::bytea), 4326)
        ELSE route_geometry
    END,
    distance_km = COALESCE(sqlc.narg(distance_km), distance_km),
    duration_min = COALESCE(sqlc.narg(duration_min), duration_min),
    annual_cost = COALESCE(sqlc.narg(annual_cost), annual_cost),
    annual_minutes = COALESCE(sqlc.narg(annual_minutes), annual_minutes),
    updated_at = NOW()
WHERE id = $1
RETURNING id, annual_cost, annual_minutes, created_at;

-- name: DeleteCommute :exec
DELETE FROM commutes WHERE id = $1;