-- Дни -------------------------------------------------------------------

-- name: CreateDay :one
INSERT INTO wakatime_days (user_id, date, total_seconds, text)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetDayByID :one
SELECT * FROM wakatime_days WHERE id = $1;

-- name: GetDayByDate :one
SELECT * FROM wakatime_days WHERE user_id = $1 AND date = $2;

-- name: ListDaysByUser :many
SELECT * FROM wakatime_days WHERE user_id = $1 ORDER BY date DESC;

-- name: UpdateDay :one
UPDATE wakatime_days
SET total_seconds = $2, text = $3, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteDay :exec
DELETE FROM wakatime_days WHERE id = $1;

-- Проекты -------------------------------------------------------------------

-- name: CreateProject :one
INSERT INTO wakatime_projects (day_id, name, total_seconds, percent, text)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetProjectByID :one
SELECT * FROM wakatime_projects WHERE id = $1;

-- name: ListProjectsByDay :many
SELECT * FROM wakatime_projects WHERE day_id = $1 ORDER BY total_seconds DESC;

-- name: UpdateProject :one
UPDATE wakatime_projects
SET total_seconds = $2, percent = $3, text = $4
WHERE id = $1
RETURNING *;

-- name: DeleteProject :exec
DELETE FROM wakatime_projects WHERE id = $1;

-- name: DeleteProjectsByDay :exec
DELETE FROM wakatime_projects WHERE day_id = $1;

-- Языки -------------------------------------------------------------------

-- name: CreateLanguage :one
INSERT INTO wakatime_languages (day_id, name, total_seconds, percent, text)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListLanguagesByDay :many
SELECT * FROM wakatime_languages WHERE day_id = $1 ORDER BY total_seconds DESC;

-- name: DeleteLanguage :exec
DELETE FROM wakatime_languages WHERE id = $1;

-- name: DeleteLanguagesByDay :exec
DELETE FROM wakatime_languages WHERE day_id = $1;

-- Редакторы -------------------------------------------------------------------

-- name: CreateEditor :one
INSERT INTO wakatime_editors (day_id, name, total_seconds, percent, text)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListEditorsByDay :many
SELECT * FROM wakatime_editors WHERE day_id = $1 ORDER BY total_seconds DESC;

-- name: DeleteEditor :exec
DELETE FROM wakatime_editors WHERE id = $1;

-- name: DeleteEditorsByDay :exec
DELETE FROM wakatime_editors WHERE day_id = $1;

-- ОС -------------------------------------------------------------------

-- name: CreateOS :one
INSERT INTO wakatime_os (day_id, name, total_seconds, percent, text)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListOSByDay :many
SELECT * FROM wakatime_os WHERE day_id = $1 ORDER BY total_seconds DESC;

-- name: DeleteOS :exec
DELETE FROM wakatime_os WHERE id = $1;

-- name: DeleteOSByDay :exec
DELETE FROM wakatime_os WHERE day_id = $1;

-- Зависимости -------------------------------------------------------------------

-- name: CreateDependency :one
INSERT INTO wakatime_dependencies (day_id, name, total_seconds, percent, text)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListDependenciesByDay :many
SELECT * FROM wakatime_dependencies WHERE day_id = $1 ORDER BY total_seconds DESC;

-- name: DeleteDependency :exec
DELETE FROM wakatime_dependencies WHERE id = $1;

-- name: DeleteDependenciesByDay :exec
DELETE FROM wakatime_dependencies WHERE day_id = $1;

-- Машины -------------------------------------------------------------------

-- name: CreateMachine :one
INSERT INTO wakatime_machines (day_id, name, total_seconds, percent, text)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListMachinesByDay :many
SELECT * FROM wakatime_machines WHERE day_id = $1 ORDER BY total_seconds DESC;

-- name: DeleteMachine :exec
DELETE FROM wakatime_machines WHERE id = $1;

-- name: DeleteMachinesByDay :exec
DELETE FROM wakatime_machines WHERE day_id = $1;

-- Сводная статистика -------------------------------------------------------------------

-- name: CreateSummary :one
INSERT INTO wakatime_summaries (user_id, start_time, end_time, range, total_seconds, daily_average, best_day_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetSummaryByID :one
SELECT * FROM wakatime_summaries WHERE id = $1;

-- name: ListSummariesByUser :many
SELECT * FROM wakatime_summaries WHERE user_id = $1 ORDER BY start_time DESC;

-- name: DeleteSummary :exec
DELETE FROM wakatime_summaries WHERE id = $1;

-- name: GetDaysByDateRange :many
SELECT * FROM wakatime_days
WHERE user_id = $1 AND date BETWEEN $2 AND $3
ORDER BY date DESC;

-- name: GetWakatimeStatsByDateRange :many
SELECT
    d.id as day_id,
    d.date,
    d.total_seconds,
    d.text as day_text,
    p.name as project_name,
    p.total_seconds as project_seconds,
    l.name as language_name,
    l.total_seconds as language_seconds
FROM
    wakatime_days d
        LEFT JOIN
    wakatime_projects p ON d.id = p.day_id
        LEFT JOIN
    wakatime_languages l ON d.id = l.day_id
WHERE
    d.user_id = $1
  AND d.date >= $2
  AND d.date <= $3
ORDER BY
    d.date DESC, p.total_seconds DESC, l.total_seconds DESC;

-- name: GetTopLanguagesByDateRange :many
SELECT
    l.name,
    SUM(l.total_seconds) as total_seconds
FROM
    wakatime_days d
    INNER JOIN
    wakatime_languages l ON d.id = l.day_id
WHERE
    d.user_id = $1
  AND d.date >= $2
  AND d.date <= $3
GROUP BY
    l.name
ORDER BY
    SUM(l.total_seconds) DESC
LIMIT $4;

-- name: GetTopProjectsByDateRange :many
SELECT
    p.name,
    SUM(p.total_seconds) as total_seconds
FROM
    wakatime_days d
    INNER JOIN
    wakatime_projects p ON d.id = p.day_id
WHERE
    d.user_id = $1
  AND d.date >= $2
  AND d.date <= $3
GROUP BY
    p.name
ORDER BY
    SUM(p.total_seconds) DESC
LIMIT $4;