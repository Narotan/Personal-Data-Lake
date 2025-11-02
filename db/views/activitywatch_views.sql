-- Представление для удобного анализа активности с человекочитаемым временем

CREATE OR REPLACE VIEW activity_summary AS
SELECT
    app,
    COUNT(*) as event_count,
    ROUND((SUM(duration) / 60)::numeric, 2) as total_minutes,
    ROUND((SUM(duration) / 3600)::numeric, 2) as total_hours,
    MIN(timestamp) as first_seen,
    MAX(timestamp) as last_seen,
    CASE
        WHEN app = 'afk' THEN '🛌 AFK'
        WHEN app = 'chrome' OR app LIKE '%browser%' THEN '🌐 Browser'
        WHEN app LIKE '%goland%' THEN '💻 IDE'
        WHEN app = 'system' THEN '⚡ System'
        ELSE '📱 ' || app
    END as category
FROM activity_events
GROUP BY app;

-- Представление для дневной статистики
CREATE OR REPLACE VIEW daily_activity_summary AS
SELECT
    date_trunc('day', timestamp)::date as day,
    app,
    COUNT(*) as event_count,
    ROUND((SUM(duration) / 60)::numeric, 2) as total_minutes,
    ROUND((SUM(duration) / 3600)::numeric, 2) as total_hours
FROM activity_events
GROUP BY day, app
ORDER BY day DESC, total_minutes DESC;

-- Представление для почасовой статистики
CREATE OR REPLACE VIEW hourly_activity AS
SELECT
    date_trunc('hour', timestamp) as hour,
    app,
    COUNT(*) as events,
    ROUND((SUM(duration) / 60)::numeric, 2) as minutes
FROM activity_events
GROUP BY hour, app
ORDER BY hour DESC;

-- Топ окон/файлов по времени
CREATE OR REPLACE VIEW top_windows AS
SELECT
    app,
    title,
    COUNT(*) as times_opened,
    ROUND((SUM(duration) / 60)::numeric, 2) as total_minutes,
    MAX(timestamp) as last_used
FROM activity_events
WHERE title IS NOT NULL AND title != '' AND title NOT LIKE '%auto-classified%'
GROUP BY app, title
ORDER BY total_minutes DESC
LIMIT 100;

COMMENT ON VIEW activity_summary IS 'Общая статистика по приложениям с категориями';
COMMENT ON VIEW daily_activity_summary IS 'Статистика активности по дням';
COMMENT ON VIEW hourly_activity IS 'Почасовая статистика активности';
COMMENT ON VIEW top_windows IS 'Топ 100 окон/файлов по времени использования';

