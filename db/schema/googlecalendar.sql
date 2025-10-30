-- Google Calendar events
CREATE TABLE IF NOT EXISTS googlecalendar_events (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    event_id VARCHAR(255) NOT NULL,
    calendar_id VARCHAR(255) NOT NULL,
    summary TEXT,
    description TEXT,
    location TEXT,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    is_all_day BOOLEAN DEFAULT FALSE,
    duration INT DEFAULT 0, -- длительность в минутах
    status VARCHAR(50) DEFAULT 'confirmed',
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    CONSTRAINT googlecalendar_events_unique UNIQUE(user_id, event_id)
);

-- Индексы для производительности
CREATE INDEX IF NOT EXISTS idx_googlecalendar_events_user_time ON googlecalendar_events(user_id, start_time DESC);
CREATE INDEX IF NOT EXISTS idx_googlecalendar_events_calendar ON googlecalendar_events(calendar_id);
CREATE INDEX IF NOT EXISTS idx_googlecalendar_events_start_time ON googlecalendar_events(start_time);

