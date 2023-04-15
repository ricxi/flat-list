#!/bin/bash
# If container keeps stopping, delete the image and try again.

cd "$(dirname "$0")"

container_name="$1"
[ -z "$container_name" ] &&
echo "must provide name for mongo container" && 
exit 1

container_status="$( docker inspect --format "{{.State.Status}}" "$container_name" )"
[ "$container_status" = "running" ] || 
[ "$container_status" = "created" ] || 
[ "$container_status" = "exited" ] &&
echo "problem trying to run container" &&
exit 1

# a throw-away dev script for starting a mongo container
docker run -d \
    --name "$container_name" \
    --publish 127.0.0.1:37017:27017 \
    --env-file usermongodb.env \
    mongo:6.0.4-jammy

# docker exec -it mongo-instance mongosh -u root -p password