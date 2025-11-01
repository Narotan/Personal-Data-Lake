-- name: InsertActivityEvent :exec
INSERT INTO activity_events (id, timestamp, duration, app, title, bucket_id)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (id) DO NOTHING;

-- name: BatchInsertActivityEvents :copyfrom
INSERT INTO activity_events (id, timestamp, duration, app, title, bucket_id)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetEventByID :one
SELECT * FROM activity_events WHERE id = $1;

-- name: GetEventsByTimeRange :many
SELECT * FROM activity_events
WHERE timestamp >= $1 AND timestamp < $2
ORDER BY timestamp DESC;

-- name: GetEventsByApp :many
SELECT * FROM activity_events
WHERE app = $1 AND timestamp >= $2 AND timestamp < $3
ORDER BY timestamp DESC;

-- name: GetTotalTimeByApp :one
SELECT COALESCE(SUM(duration), 0) as total_seconds
FROM activity_events
WHERE app = $1 AND timestamp >= $2 AND timestamp < $3;

-- name: GetAppStats :many
SELECT
    app,
    COUNT(*) as event_count,
    SUM(duration) as total_seconds
FROM activity_events
WHERE timestamp >= $1 AND timestamp < $2
GROUP BY app
ORDER BY total_seconds DESC;

-- name: GetRecentEvents :many
SELECT * FROM activity_events
ORDER BY timestamp DESC
LIMIT $1;

-- name: DeleteOldEvents :exec
DELETE FROM activity_events
WHERE timestamp < $1;

