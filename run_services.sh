#!/bin/bash
# This is a development script that I use to quickly spin up the app locally for some quick tests.

# called to load all the environment variables from
# a specific .env file for a go service before running it
load_env_file() {
    env_file="$1"
    export $(grep -v '^#' "$env_file" | xargs)
}

cleanup() {
    echo -e "\ncleaning up..."
    sleep 5
    # kill -9 "$pid"
    exit 0
}

trap cleanup SIGINT

(
    # start the token service
    load_env_file token.env &&
    ./dev_scripts/db/run_tokendb.sh tokendb-instance 2>>errors.txt &&
    ./dev_scripts/db/run_tokendb_migrations.sh tokendb-instance 2>>errors.txt
    cd token &&
    go run ./cmd/grpc
	# go run ./cmd/http
) &
(
    # start the mailer service
    load_env_file mailer.env &&
    cd mailer &&
    go run ./cmd/grpc
    # go run ./cmd/http
) &
(
    # start react mailer client
	cd frontend-client && {
        [ ! -d node_modules ] && npm i 
        npm run dev 
    } 
) &
(
    # start the user service
    load_env_file user.env &&
    ./dev_scripts/db/run_flatlistdb.sh flatlistdb-instance 2>>errors.txt
    cd user &&
    go run ./cmd/http
) &
(
    # start the task service
    load_env_file task.env &&
    # this depends on the same database as the user service, which should already be running
    sleep 5 # I'll write a more sustainable solution later than sleeping
    cd task &&
    go run ./cmd/http
    # debug:
    # ./taskservice
    # go build -gcflags=all="-N -l" -o taskservice ./cmd/http
) &
(
    # list running services
    # update this so I can input the ports dynamically
    # while true; do
    sleep 10
    lsof -i :5000-5009 -i :9000 | awk '{print $1, $2, $5, $8, $9}'
    # done
) &
# pid=$!

# do not remove wait because:
# it allows time for our go services to clean up and gracefully shutdown
# it is also needed if I want to run a trap
wait