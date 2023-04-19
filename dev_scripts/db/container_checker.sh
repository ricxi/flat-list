#!/bin/bash
# utility functions for checking the status of docker containers

# can receive either the container name or id as an argument
# call this script to check if a container with a specific name or id exists.
# If the container exists, then exit 0; otherwise, exit 1.
if_exists() {
    [ "$#" -eq 0 ] &&
    echo "Usage: $0 <container_name_or_id>" &&
    return 1

    local container=$1

    local status="$( docker inspect --format "{{.State.Status}}" "$container" )"
    [ "$?" -ne 0 ] && 
    return 1 # catch the error from the docker inspect command and exit

    [ "$status" = "running" ] ||
    [ "$status" = "created" ] || 
    [ "$status" = "exited" ] &&
    return 0

    return 1
}

is_running() {
    [ "$#" -eq 0 ] &&
    echo "Usage: $0 <container_name_or_id>" &&
    return 1

    local container=$1

    local status="$( docker inspect --format "{{.State.Status}}" "$container" )"
    [ "$?" -ne 0 ] && 
    return 1 # catch the error from the docker inspect command and return exit code 1

    [ "$status" = "running" ] &&
    return 0

    return 1
}