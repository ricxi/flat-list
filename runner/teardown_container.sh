#!/bin/sh
# THIS IS A SHELL SCRIPT, NOT A BASH SCRIPT

# This can be either the name or id of the container
container=$1
containerStatus="$( docker inspect --format "{{.State.Status}}" "$container" )"

# If the container is running, stop it and remove it.
# If its already been stopped (status=exited), then just remove it
if [ "$containerStatus" = "running" ] || [ "$containerStatus" = "exited" ]; then
    # This line assumes that the container is succesfully stopped before we delete it.
    # So if the container isn't stopped, then there's going to be some weird error messages.
    [ "$containerStatus" = "running" ] && [ "$( docker stop "$container" )" ] 
    docker rm "$container"
fi

# one line for stopping and removing the container
# containerName=$1 && docker rm $(docker stop $containerName 2>/dev/null) 2>/dev/null