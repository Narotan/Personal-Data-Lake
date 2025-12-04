#!/bin/bash
set -e

echo "🗃️  Initializing database schema..."

# Run all migration files in order
for migration in /docker-entrypoint-initdb.d/migrations/*.up.sql; do
    if [ -f "$migration" ]; then
        echo "Running migration: $(basename $migration)"
        psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" -f "$migration"
    fi
done

# Create default user if API_USER_ID is set
if [ ! -z "$API_USER_ID" ]; then
    echo "Creating default user with ID: $API_USER_ID"
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
        INSERT INTO users (id, username, email)
        VALUES ('$API_USER_ID', 'default_user', 'user@example.com')
        ON CONFLICT (id) DO NOTHING;
EOSQL
    echo "✅ Default user created or already exists"
fi

echo "✅ Database initialization complete!"

