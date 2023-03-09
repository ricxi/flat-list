# import environment variables
include .envrc

USER_BINARY=userService

# run outside of docker container
.PHONY: run/dev/user
run/dev/user: 
	@echo "DEV: starting user service on port ${PORT}..."
	cd ./user && go run ./cmd/http/

.PHONY: build/user
build/user:
	@echo "building binary..."
	cd ./user && go build -o bin/${USER_BINARY} ./cmd/http/

# start a mongo container for the user service
.PHONY: run/dev/mongo
run/dev/mongo:
	@echo "DEV: running local mongo container..."
	cd scripts && ./mongo.sh