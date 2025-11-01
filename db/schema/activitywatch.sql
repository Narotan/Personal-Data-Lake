CREATE TABLE IF NOT EXISTS activity_events (
    id BIGINT PRIMARY KEY,
    timestamp TIMESTAMPTZ NOT NULL,
    duration FLOAT NOT NULL,
    app TEXT NOT NULL,
    title TEXT,
    bucket_id TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_activity_timestamp ON activity_events(timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_activity_app ON activity_events(app);
CREATE INDEX IF NOT EXISTS idx_activity_bucket ON activity_events(bucket_id);

