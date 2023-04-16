#!/bin/bash
# another throw away script
cd "$(dirname "$0")"

source ./container_checker.sh

container_name="$1"
[ -z "$container_name" ] &&
echo "must provide name for postgres container" && 
exit 1

if_exists "$container_name" &&
is_running "$container_name" && 
exit 0

if_exists "$container_name" || 
is_running "$container_name" &&
exit 1 # docker start "$container_name"

# spin up a postgres container 
# if container keeps stopping, delete the image and try again.
docker run -d \
    --name "$container_name" \
    --publish 127.0.0.1:5433:5432 \
    --env-file=tokenpsqldb.env \
    postgres:15.2-alpine
