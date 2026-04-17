-- name: CreateURL :one
INSERT INTO urls (
    short_code,
    original_url,
    user_id,
    expires_at,
    max_clicks
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING id, short_code, original_url, user_id, is_active, expires_at, max_clicks, created_at, updated_at;


-- name: GetURLByShortCode :one
SELECT id, short_code, original_url, user_id, is_active, expires_at, max_clicks, created_at, updated_at
FROM urls
WHERE short_code = $1
LIMIT 1;


-- name: GetURLsByUserID :many
SELECT id, short_code, original_url, user_id, is_active, expires_at, max_clicks, created_at, updated_at
FROM urls
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;


-- name: DeactivateURL :exec
UPDATE urls
SET is_active = false,
    updated_at = NOW()
WHERE id = $1;


-- name: UpdateURL :one
UPDATE urls
SET original_url = $2,
    expires_at    = $3,
    max_clicks    = $4,
    updated_at    = NOW()
WHERE id = $1
RETURNING id, short_code, original_url, user_id, is_active, expires_at, max_clicks, created_at, updated_at;


-- name: DeleteExpiredURLs :execrows
DELETE FROM urls
WHERE expires_at IS NOT NULL
AND expires_at < NOW();