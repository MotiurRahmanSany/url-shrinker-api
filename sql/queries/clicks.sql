-- name: CreateClick :one
INSERT INTO clicks (
    url_id,
    ip_address,
    user_agent,
    referer
) VALUES (
    $1, $2, $3, $4
)
RETURNING id, url_id, clicked_at, ip_address, user_agent, referer;


-- name: CountClicksByURLID :one
SELECT COUNT(*) AS total
FROM clicks
WHERE url_id = $1;


-- name: CountClicksTodayByURLID :one
SELECT COUNT(*) AS total
FROM clicks
WHERE url_id = $1
  AND clicked_at >= CURRENT_DATE;


-- name: GetClicksByURLIDGroupedByDay :many
SELECT
    DATE(clicked_at) AS day,
    COUNT(*)         AS total
FROM clicks
WHERE url_id = $1
GROUP BY DATE(clicked_at)
ORDER BY day ASC;
