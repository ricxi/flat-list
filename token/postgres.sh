#!/bin/bash

# get dsn for postgres connection
psql_dsn=$1
if [ -z "$psql_dsn" ]; then
    echo "please enter postgres connection string as command-line argument"
    exit 1
fi

# spin up a postgres container 
docker run -d \
    --name postgres-instance \
    --publish 127.0.0.1:5433:5432 \
    --env-file=.postgres.env \
    postgres:15.2-alpine

# make 3ish attempts to check if postgres is ready to accept connections
tries=0
until [ "$(docker exec postgres-instance pg_isready | grep "accepting connections" )" ] || [ "$tries" -gt 3 ]; do
    echo "waiting for postgres to accept connections"
    sleep 2
    tries=$((tries+1))
done

# run migrations if postgres is ready to accept connections
if [ "$tries" -le 3 ]; then
    migrate -path=./migrations -database "$psql_dsn" up
fi
