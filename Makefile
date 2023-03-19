# import environment variables
include .envrc

USER_BINARY=userService

# run outside of docker container
.PHONY: run/dev/user
run/dev/user: 
	@echo "DEV: starting user service on port ${PORT}..."
	cd ./user && go run ./cmd/http/

# test service layer of user service
.PHONY: test/user/service
test/user/service:
	@echo "TEST: user service layer"
	cd ./user && go test -v mocks_test.go service_test.go

.PHONY: build/user
build/user:
	@echo "building binary..."
	cd ./user && go build -o bin/${USER_BINARY} ./cmd/http/

# start a mongo container for the user service
.PHONY: run/dev/mongo
run/dev/mongo:
	@echo "DEV: running local mongo container..."
	cd dev_scripts && ./start_mongo.sh

# start front-end react client for email activation 
.PHONY: run/dev/react/email
run/dev/react/email:
	@echo "DEV: running react email client"
	cd frontend-client && npm run dev