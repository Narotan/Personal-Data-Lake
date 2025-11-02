-- name: BulkInsertEvents :copyfrom
INSERT INTO activity_events (
    timestamp,
    duration,
    app,
    title,
    bucket_id
) VALUES (
    $1, $2, $3, $4, $5
);

-- name: GetAppStats :many
SELECT
    app,
    SUM(duration)::float as total_duration,
    COUNT(*) as event_count
FROM activity_events
WHERE timestamp >= $1 AND timestamp < $2
GROUP BY app
ORDER BY total_duration DESC;

-- name: GetRecentEvents :many
SELECT * FROM activity_events
WHERE timestamp >= $1
ORDER BY timestamp DESC
LIMIT $2;

-- name: GetEventsByApp :many
SELECT * FROM activity_events
WHERE app = $1 AND timestamp >= $2 AND timestamp < $3
ORDER BY timestamp DESC;

