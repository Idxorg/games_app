#!/bin/sh
set -e

# Migration runner — waits for PostgreSQL then applies .sql files in order.
#
# Required environment variables:
#   DATABASE_URL — full PostgreSQL connection string
#                  e.g. postgresql://game:pass@db:5432/game_platform?sslmode=disable
#
# The script extracts host, port, user, and database from DATABASE_URL to run
# pg_isready for the readiness check, then feeds each .sql file to psql.

if [ -z "$DATABASE_URL" ]; then
  echo "ERROR: DATABASE_URL environment variable is not set"
  exit 1
fi

# Parse DATABASE_URL: postgresql://user:password@host:port/dbname?params
# Extract components using parameter expansion
DB_PROTO="${DATABASE_URL%%://*}"
DB_NO_PROTO="${DATABASE_URL#*://}"
DB_CREDS="${DB_NO_PROTO%%@*}"
DB_USER="${DB_CREDS%%:*}"
DB_REST="${DB_NO_PROTO#*@}"
DB_HOST_PORT="${DB_REST%%/*}"
DB_HOST="${DB_HOST_PORT%%:*}"
DB_PORT="${DB_HOST_PORT#*:}"
DB_PATH="${DB_REST#*/}"
DB_NAME="${DB_PATH%%\?*}"

echo "Waiting for PostgreSQL at ${DB_HOST}:${DB_PORT} ..."

# Wait up to 60 seconds for PostgreSQL to become ready
RETRIES=30
until pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -q 2>/dev/null; do
  RETRIES=$((RETRIES - 1))
  if [ "$RETRIES" -le 0 ]; then
    echo "ERROR: PostgreSQL did not become ready in time"
    exit 1
  fi
  sleep 2
done

echo "PostgreSQL is ready. Running migrations ..."

MIGRATION_DIR="$(dirname "$0")"

# Run each .sql file in sorted order
for SQL_FILE in "$MIGRATION_DIR"/*.sql; do
  [ -f "$SQL_FILE" ] || continue
  echo "Running migration: $(basename "$SQL_FILE")"
  psql "$DATABASE_URL" -f "$SQL_FILE" --quiet --set ON_ERROR_STOP=1
  if [ $? -ne 0 ]; then
    echo "ERROR: Migration failed: $(basename "$SQL_FILE")"
    exit 1
  fi
done

echo "All migrations completed successfully."
