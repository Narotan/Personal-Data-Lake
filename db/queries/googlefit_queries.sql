-- Daily Stats Queries -------------------------------------------------------------------

-- name: CreateDailyStat :one
INSERT INTO googlefit_daily_stats (user_id, date, steps, distance)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetDailyStatByDate :one
SELECT * FROM googlefit_daily_stats WHERE user_id = $1 AND date = $2;

-- name: ListDailyStatsByUser :many
SELECT * FROM googlefit_daily_stats
WHERE user_id = $1
ORDER BY date DESC
LIMIT $2 OFFSET $3;

-- name: ListDailyStatsByDateRange :many
SELECT * FROM googlefit_daily_stats
WHERE user_id = $1 AND date BETWEEN $2 AND $3
ORDER BY date DESC;

-- name: UpdateDailyStat :one
UPDATE googlefit_daily_stats
SET steps = $2, distance = $3, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteDailyStat :exec
DELETE FROM googlefit_daily_stats WHERE id = $1;

-- name: UpsertDailyStat :one
INSERT INTO googlefit_daily_stats (user_id, date, steps, distance)
VALUES ($1, $2, $3, $4)
ON CONFLICT (user_id, date)
DO UPDATE SET
    steps = EXCLUDED.steps,
    distance = EXCLUDED.distance,
    updated_at = now()
RETURNING *;

-- Analytics Queries -------------------------------------------------------------------

-- name: GetWeeklyStepsSummary :many
SELECT
    date,
    steps,
    distance,
    EXTRACT(DOW FROM date) as day_of_week
FROM googlefit_daily_stats
WHERE user_id = $1
  AND date >= CURRENT_DATE - INTERVAL '7 days'
ORDER BY date DESC;

-- name: GetMonthlyAverage :one
SELECT
    AVG(steps) as avg_steps,
    AVG(distance) as avg_distance,
    SUM(steps) as total_steps,
    SUM(distance) as total_distance
FROM googlefit_daily_stats
WHERE user_id = $1
  AND date >= CURRENT_DATE - INTERVAL '30 days';


-- name: GetGoogleFitDailyStatsByDateRange :many
SELECT
    date,
    steps,
    distance
FROM
    googlefit_daily_stats
WHERE
    user_id = $1
  AND date >= $2
  AND date <= $3
ORDER BY
    date DESC;