#!/bin/bash

cd "$(dirname "$0")"

# spin up a postgres container 
docker run -d \
    --name postgres-instance \
    --publish 127.0.0.1:5433:5432 \
    --env-file=tokenpsqldb.env \
    postgres:15.2-alpine
