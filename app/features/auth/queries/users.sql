-- name: GetUser :one
SELECT * FROM users WHERE id = ? LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY username;

-- name: CreateUser :one
INSERT INTO users (id, username, email, password_hash) 
VALUES (?, ?, ?, ?) 
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = ? LIMIT 1;

-- name: UpdateUser :exec
UPDATE users SET username = ?, email = ? WHERE id = ?;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?;

-- name: CheckIfUserExistsByEmail :one
SELECT EXISTS(SELECT 1 FROM users WHERE email = ?);

-- name: CheckIfUserExistsByUsername :one
SELECT EXISTS(SELECT 1 FROM users WHERE username = ?);