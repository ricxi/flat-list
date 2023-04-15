#!/bin/bash

# called to load all the environment variables from
# a specific .env file for a go service before running it
load_env_file() {
    env_file="$1"
    export $(grep -v '^#' "$env_file" | xargs)
}

cleanup() {
    echo -e "\ncleaning up..."
}

trap cleanup SIGINT

(
    # start the token service
    load_env_file token.env &&
    ./dev_scripts/db/run_tokendb.sh &&
    ./dev_scripts/db/run_tokendb_migrations.sh &&
    cd token &&
    go run ./cmd/grpc
	# go run ./cmd/http
) &
echo "$!"

(
    # start the mailer service
    load_env_file mailer.env &&
    cd mailer &&
    go run ./cmd/grpc
    # go run ./cmd/http
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
    ./dev_scripts/db/run_userdb.sh &&
    cd user &&
    go run ./cmd/http
) &
echo "$!"

# do not remove wait because:
# it allows time for our go services to clean up and gracefully shutdown
# it is also needed if I want to run a trap
wait