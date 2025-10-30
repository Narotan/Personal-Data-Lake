-- Google Fit daily statistics
CREATE TABLE IF NOT EXISTS googlefit_daily_stats (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    steps INT DEFAULT 0,
    distance FLOAT DEFAULT 0, -- в метрах
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    CONSTRAINT googlefit_daily_stats_unique UNIQUE(user_id, date)
);

-- Индекс для производительности
CREATE INDEX IF NOT EXISTS idx_googlefit_daily_stats_user_date ON googlefit_daily_stats(user_id, date DESC);

