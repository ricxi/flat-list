#!/bin/bash
# Run migrations on postgres container 

# Refer to 'start_postgres.sh' file if PSQL_DSN doesn't work properly
PSQL_DSN="postgresql://postgres:password@127.0.0.1:5433/tokens?sslmode=disable"

# Make 3ish attempts to check if postgres is ready to accept connections
tries=0
until [ "$(docker exec postgres-instance pg_isready | grep "accepting connections")" ]; do
    [ "$tries" -ge 3 ] && exit 1
    echo "waiting for postgres to accept connections"
    tries=$((tries+1))
    sleep "$((tries+2))"
done

# Run migrations when postgres is ready to accept connections
migrate -path=../migrations/token/ -database "$PSQL_DSN" up