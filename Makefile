# Makefiles in this application are used mainly for testing
include .envrc

USER_BINARY=userService

# test service layer of user service
.PHONY: test/user/service
test/user/service:
	@echo "TEST: user service layer"
	cd ./user && go test -v -run Test_Service

.PHONY: test/e2e/user
test/e2e/user:
	@echo "TEST E2E: user microservice running..."
	cd ./user/cmd/http && go test -v

.PHONY: build/user
build/user:
	@echo "building binary..."
	cd ./user && go build -o bin/${USER_BINARY} ./cmd/http/

# test for task service
.PHONY: test/task
test/task:
	@echo "TEST: task microservice running..."
	cd ./task && go test -cover

.PHONY: test/e2e/task
test/e2e/task:
	@echo "TEST E2E: task microservice running..."
	cd ./task && go test -v ./cmd/http