-- name: GetCommute :one
SELECT * FROM commutes WHERE id = $1;

-- name: CreateCommute :one
INSERT INTO commutes (
    user_id, name, home_point, office_point, route_geometry,
    distance_km, duration_min, vehicle, fuel_price, days_per_week,
    annual_cost, annual_minutes
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
)
RETURNING id, annual_cost, annual_minutes, created_at;

-- name: ListCommutesByUser :many
SELECT *
FROM commutes WHERE user_id = $1 ORDER BY created_at DESC;

-- name: UpdateCommute :one
UPDATE commutes SET name = $2, updated_at = NOW() WHERE id = $1 RETURNING id;

-- name: DeleteCommute :exec
DELETE FROM commutes WHERE id = $1;
