# Makefiles in this application are used mainly for testing
include .envrc

USER_BINARY=userService

# test service layer of user service
.PHONY: test/user/service
test/user/service:
	@echo "TEST: user service layer"
	cd ./user && go test -v -run Test_Service

.PHONY: test/user/e2e
test/user/e2e:
	@echo "TEST E2E: user microservice"
	cd ./user/cmd/http && go test -v

.PHONY: build/user
build/user:
	@echo "building binary..."
	cd ./user && go build -o bin/${USER_BINARY} ./cmd/http/

.PHONY: tidy/user
tidy/user:
	@echo "Tidying up user service"
	cd ./user && go mod tidy