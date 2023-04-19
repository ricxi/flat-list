#!/bin/bash
# a throw-away dev script for starting a mongo container
cd "$(dirname "$0")"

source ./container_checker.sh

container_name="$1"
[ -z "$container_name" ] &&
echo "must provide name for mongo container" && 
exit 1

# If the container exists and it's running, then exit 0.
# If the container does not exist and it's not running, then continue to docker run.
if_exists "$container_name" &&
is_running "$container_name" && 
exit 0


# It's impossible for a container that doesn't exist to run,
# so the only possibility here is if a container exists, and it's not running,
# then we exit 1 (instead of using 'docker start' because I don't know the state of the container) 
if_exists "$container_name" || 
is_running "$container_name" &&
exit 1 # docker start "$container_name"

# -v ./mongo-init.sh:/docker-entrypoint-initdb.d/mongo-init.sh:ro \
docker run -d \
    --name "$container_name" \
    --publish 127.0.0.1:37017:27017 \
    --env-file mongodb.env \
    mongo:6.0.4-jammy

# docker exec -it userdb-instance mongosh -u root -p password