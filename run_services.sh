#!/bin/bash

# called to load all the environment variables from
# a specific .env file for a go service before running it
load_env_file() {
    env_file="$1"
    export $(grep -v '^#' "$env_file" | xargs)
}

cleanup() {
    echo -e "\ncleaning up"
}

trap cleanup SIGINT

(
    # start the token service
    load_env_file token.env &&
    # ./dev_scripts/start_postgres.sh postgresql://postgres:password@127.0.0.1:5433/tokens?sslmode=disable &&
    cd token &&
    go run ./cmd/grpc
) &
echo "$!"

(
    # start the mailer service
    load_env_file mailer.env &&
    cd mailer &&
    go run ./cmd/grpc
) &
echo "$!"

(
    # start react mailer client
	cd frontend-client && npm run dev
) &
echo "$!"

(
    # start the user service
    load_env_file user.env &&
    # ./dev_scripts/start_mongo.sh &&
    cd user &&
    go run ./cmd/http
) &
echo "$!"

# do not remove wait:
# allows time for our go services to clean up and gracefully shutdown
# also needed if I want to run a trap
wait