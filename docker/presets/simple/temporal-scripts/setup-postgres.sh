#!/bin/sh
set -eu

echo 'Starting PostgreSQL schema setup for Temporal...'
echo 'Waiting for PostgreSQL port to be available...'

# Wait for postgres to be reachable
max_attempts=30
attempt=0
until nc -z -w 5 ${POSTGRES_SEEDS} ${DB_PORT}; do
  attempt=$((attempt + 1))
  if [ $attempt -ge $max_attempts ]; then
    echo "PostgreSQL did not become available after $max_attempts attempts"
    exit 1
  fi
  echo "Waiting for PostgreSQL... (attempt $attempt/$max_attempts)"
  sleep 2
done
echo 'PostgreSQL port is available'

# Create and setup temporal database
echo 'Creating temporal database...'
temporal-sql-tool --plugin postgres12 --ep ${POSTGRES_SEEDS} -u ${POSTGRES_USER} -p ${DB_PORT} --db temporal create || echo 'Database temporal may already exist'
temporal-sql-tool --plugin postgres12 --ep ${POSTGRES_SEEDS} -u ${POSTGRES_USER} -p ${DB_PORT} --db temporal setup-schema -v 0.0
temporal-sql-tool --plugin postgres12 --ep ${POSTGRES_SEEDS} -u ${POSTGRES_USER} -p ${DB_PORT} --db temporal update-schema -d /etc/temporal/schema/postgresql/v12/temporal/versioned

# Create and setup visibility database
echo 'Creating temporal_visibility database...'
temporal-sql-tool --plugin postgres12 --ep ${POSTGRES_SEEDS} -u ${POSTGRES_USER} -p ${DB_PORT} --db temporal_visibility create || echo 'Database temporal_visibility may already exist'
temporal-sql-tool --plugin postgres12 --ep ${POSTGRES_SEEDS} -u ${POSTGRES_USER} -p ${DB_PORT} --db temporal_visibility setup-schema -v 0.0
temporal-sql-tool --plugin postgres12 --ep ${POSTGRES_SEEDS} -u ${POSTGRES_USER} -p ${DB_PORT} --db temporal_visibility update-schema -d /etc/temporal/schema/postgresql/v12/visibility/versioned

echo 'PostgreSQL schema setup complete'
