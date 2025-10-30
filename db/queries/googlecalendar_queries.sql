-- Events Queries -------------------------------------------------------------------

-- name: CreateEvent :one
INSERT INTO googlecalendar_events (
    user_id, event_id, calendar_id, summary, description, location,
    start_time, end_time, is_all_day, duration, status
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: GetEventByID :one
SELECT * FROM googlecalendar_events WHERE user_id = $1 AND event_id = $2;

-- name: ListEventsByUser :many
SELECT * FROM googlecalendar_events
WHERE user_id = $1
ORDER BY start_time DESC
LIMIT $2 OFFSET $3;

-- name: ListEventsByDateRange :many
SELECT * FROM googlecalendar_events
WHERE user_id = $1 AND start_time >= $2 AND end_time <= $3
ORDER BY start_time ASC;

-- name: ListEventsByDate :many
SELECT * FROM googlecalendar_events
WHERE user_id = $1 AND DATE(start_time) = $2
ORDER BY start_time ASC;

-- name: ListEventsByCalendar :many
SELECT * FROM googlecalendar_events
WHERE user_id = $1 AND calendar_id = $2
ORDER BY start_time DESC
LIMIT $3 OFFSET $4;

-- name: UpdateEvent :one
UPDATE googlecalendar_events
SET summary = $2, description = $3, location = $4,
    start_time = $5, end_time = $6, is_all_day = $7,
    duration = $8, status = $9, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteEvent :exec
DELETE FROM googlecalendar_events WHERE id = $1;

-- name: UpsertEvent :one
INSERT INTO googlecalendar_events (
    user_id, event_id, calendar_id, summary, description, location,
    start_time, end_time, is_all_day, duration, status
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
ON CONFLICT (user_id, event_id)
DO UPDATE SET
    calendar_id = EXCLUDED.calendar_id,
    summary = EXCLUDED.summary,
    description = EXCLUDED.description,
    location = EXCLUDED.location,
    start_time = EXCLUDED.start_time,
    end_time = EXCLUDED.end_time,
    is_all_day = EXCLUDED.is_all_day,
    duration = EXCLUDED.duration,
    status = EXCLUDED.status,
    updated_at = now()
RETURNING *;

-- Analytics Queries -------------------------------------------------------------------

-- name: GetDailyEventsSummary :many
SELECT
    DATE(start_time) as event_date,
    COUNT(*) as total_events,
    SUM(duration) as total_duration_minutes,
    COUNT(CASE WHEN is_all_day THEN 1 END) as all_day_events
FROM googlecalendar_events
WHERE user_id = $1
  AND start_time >= $2
  AND end_time <= $3
GROUP BY DATE(start_time)
ORDER BY event_date DESC;

-- name: GetEventsByCalendarSummary :many
SELECT
    calendar_id,
    COUNT(*) as total_events,
    SUM(duration) as total_duration_minutes
FROM googlecalendar_events
WHERE user_id = $1
  AND start_time >= $2
  AND end_time <= $3
GROUP BY calendar_id
ORDER BY total_events DESC;

-- name: GetBusiestDays :many
SELECT
    DATE(start_time) as event_date,
    COUNT(*) as event_count,
    SUM(duration) as total_duration_minutes
FROM googlecalendar_events
WHERE user_id = $1
  AND start_time >= CURRENT_DATE - INTERVAL '30 days'
GROUP BY DATE(start_time)
ORDER BY event_count DESC, total_duration_minutes DESC
LIMIT $2;

-- name: GetAverageDailyEvents :one
SELECT
    AVG(event_count)::FLOAT as avg_events_per_day,
    AVG(total_duration)::FLOAT as avg_duration_per_day
FROM (
    SELECT
        DATE(start_time) as event_date,
        COUNT(*) as event_count,
        SUM(duration) as total_duration
    FROM googlecalendar_events
    WHERE user_id = $1
      AND start_time >= $2
      AND end_time <= $3
    GROUP BY DATE(start_time)
) AS daily_stats;

