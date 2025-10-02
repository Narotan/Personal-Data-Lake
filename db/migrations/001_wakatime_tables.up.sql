-- Пользователи (если ещё нет)
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    username TEXT NOT NULL,
    email TEXT UNIQUE,
    created_at TIMESTAMPTZ DEFAULT now()
);

-- Дни
CREATE TABLE IF NOT EXISTS wakatime_days (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    total_seconds FLOAT NOT NULL,
    text TEXT,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    CONSTRAINT wakatime_days_unique UNIQUE(user_id, date)
);

-- Проекты
CREATE TABLE IF NOT EXISTS wakatime_projects (
    id SERIAL PRIMARY KEY,
    day_id INT NOT NULL REFERENCES wakatime_days(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    total_seconds FLOAT NOT NULL,
    percent FLOAT,
    text TEXT,
    created_at TIMESTAMPTZ DEFAULT now()
);

-- Языки
CREATE TABLE IF NOT EXISTS wakatime_languages (
    id SERIAL PRIMARY KEY,
    day_id INT NOT NULL REFERENCES wakatime_days(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    total_seconds FLOAT NOT NULL,
    percent FLOAT,
    text TEXT,
    created_at TIMESTAMPTZ DEFAULT now()
);

-- Редакторы
CREATE TABLE IF NOT EXISTS wakatime_editors (
    id SERIAL PRIMARY KEY,
    day_id INT NOT NULL REFERENCES wakatime_days(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    total_seconds FLOAT NOT NULL,
    percent FLOAT,
    text TEXT,
    created_at TIMESTAMPTZ DEFAULT now()
);

-- Операционные системы
CREATE TABLE IF NOT EXISTS wakatime_os (
    id SERIAL PRIMARY KEY,
    day_id INT NOT NULL REFERENCES wakatime_days(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    total_seconds FLOAT NOT NULL,
    percent FLOAT,
    text TEXT,
    created_at TIMESTAMPTZ DEFAULT now()
);

-- Зависимости
CREATE TABLE IF NOT EXISTS wakatime_dependencies (
    id SERIAL PRIMARY KEY,
    day_id INT NOT NULL REFERENCES wakatime_days(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    total_seconds FLOAT NOT NULL,
    percent FLOAT,
    text TEXT,
    created_at TIMESTAMPTZ DEFAULT now()
);

-- Машины
CREATE TABLE IF NOT EXISTS wakatime_machines (
    id SERIAL PRIMARY KEY,
    day_id INT NOT NULL REFERENCES wakatime_days(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    total_seconds FLOAT NOT NULL,
    percent FLOAT,
    text TEXT,
    created_at TIMESTAMPTZ DEFAULT now()
);

-- Сводная статистика
CREATE TABLE IF NOT EXISTS wakatime_summaries (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    range TEXT,
    total_seconds FLOAT,
    daily_average FLOAT,
    best_day_id INT REFERENCES wakatime_days(id),
    created_at TIMESTAMPTZ DEFAULT now()
);
