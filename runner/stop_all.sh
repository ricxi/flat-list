#!/bin/bash

# stop ALL docker containers
docker ps -q | xargs docker stop | xargs docker rm