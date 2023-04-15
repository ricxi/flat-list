#!/bin/bash

# If container keeps stopping, delete the image and try again.
cd "$(dirname "$0")"

# a throw-away dev script for starting a mongo container
docker run -d \
    --name mongo-instance \
    --publish 127.0.0.1:37017:27017 \
    --env-file usermongodb.env \
    mongo:6.0.4-jammy

# docker exec -it mongo-instance mongosh -u root -p password