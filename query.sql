-- name: GetBoops :many
SELECT * FROM boops ORDER BY id DESC;

-- name: GetBoopsTasks :many
SELECT * FROM boops WHERE lower(text) LIKE '%todo%' OR lower(text) LIKE '%doing%' OR lower(text) LIKE '%done%' ORDER BY id DESC;

-- name: GetBoopsFolder :many
SELECT * FROM boops WHERE starts_with(text, $1) ORDER BY id DESC;

-- name: CreateBoop :one
INSERT INTO boops (
    text
) VALUES (
    $1
)
RETURNING *;

-- name: UpdateBoop :one
UPDATE boops
SET text = $2
WHERE id = $1
RETURNING *;

-- name: DeleteBoop :exec
DELETE FROM boops
WHERE id = $1;
