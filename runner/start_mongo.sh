#!/bin/bash

# spin up a mongo container
# make sure that the 'mongo.env' file is in the same directory
# this script must always return the docker container's id
docker run -d \
    --name mongo-instance \
    --publish 127.0.0.1:37017:27017 \
    --env-file mongo.env \
    mongo:6.0.4-jammy

# docker exec -it mongo-instance mongosh -u root -p password