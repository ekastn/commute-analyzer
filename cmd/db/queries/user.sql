-- name: CreateUser :one
INSERT INTO users (device_id)
VALUES ($1)
RETURNING user_id;

-- name: GetUserByDeviceId :one
SELECT user_id FROM users WHERE device_id = $1 LIMIT 1;