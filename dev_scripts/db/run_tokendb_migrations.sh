#!/bin/bash
# Run migrations on postgres container for the token service database

# Refer to 'run_tokendb.sh' file if PSQL_DSN doesn't work properly
PSQL_DSN="postgresql://postgres:password@127.0.0.1:5433/tokens?sslmode=disable"
# Path to migrations directory for token database
MIGRATION_FILES="../../migrations/token/"

cd "$(dirname "$0")"

container_name="$1"

# Make 3ish attempts to check if postgres is ready to accept connections
tries=0
until [ "$(docker exec "$container_name" pg_isready | grep "accepting connections")" ]; do
    [ "$tries" -ge 3 ] && exit 1
    echo "waiting for postgres to accept connections"
    tries=$((tries+1))
    sleep "$((tries+2))"
done

# Run migrations when postgres is ready to accept connections
migrate -path=../../migrations/token/ -database "$PSQL_DSN" up