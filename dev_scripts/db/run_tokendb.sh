#!/bin/bash
# If container keeps stopping, delete the image and try again.

cd "$(dirname "$0")"

container_name="$1"
[ -z "$container_name" ] &&
echo "must provide name for postgres container" && 
exit 1


container_status="$( docker inspect --format '{{.State.Status}}' "$container_name" )"
[ "$container_status" = "running" ] || 
[ "$container_status" = "created" ] || 
[ "$container_status" = "exited" ] &&
echo "problem trying to run container" &&
exit 1

# spin up a postgres container 
docker run -d \
    --name "$container_name" \
    --publish 127.0.0.1:5433:5432 \
    --env-file=tokenpsqldb.env \
    postgres:15.2-alpine
