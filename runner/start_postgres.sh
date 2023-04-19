#!/bin/bash

cd "$(dirname "$0")"

# spin up a postgres container 
# make sure that the postgres.env file is in the same directory
# this script must always return the docker container's id
# If host, port, password, or whatever changes, then update the psql_migrations.sh file
docker run -d \
    --name postgres-instance \
    --publish 127.0.0.1:5433:5432 \
    --env-file=postgres.env \
    postgres:15.2-alpine