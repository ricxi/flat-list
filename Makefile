# Makefiles in this application are used mainly for testing
include .envrc

USER_BINARY=userService

# test service layer of user service
.PHONY: test/user/service
test/user/service:
	@echo "TEST: user service layer"
	cd ./user && go test -v -run Test_Service

.PHONY: test/user
test/user:
	@echo "TEST: user service"
	cd ./user && go test

.PHONY: test/e2e/user
test/e2e/user:
	@echo "TEST E2E: user microservice running..."
	cd ./user/cmd/http && go test -v

.PHONY: build/user
build/user:
	@echo "building binary..."
	cd ./user && go build -o bin/${USER_BINARY} ./cmd/http/

# test for task service
# ! requires a database connection
.PHONY: test/task
test/task:
	@echo "TEST: task microservice running..."
	cd ./task && go test -cover

# ! requires a database connection
.PHONY: test/e2e/task
test/e2e/task:
	@echo "TEST E2E: task microservice running..."
	cd ./task && go test ./cmd/http

# generate/recompile grpc code for mailer service
.PHONY: protoc/mailer
protoc/mailer:
	cd mailer && protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pb/mailer.proto

# build the lambda mailer (make sure to export .envrc)
.PHONY: build/lambda
build/lambda:
	cd mailer && GOOS=linux GOARCH=amd64 go build -o bin/lambdaMailer ./cmd/lambda && cp -r templates ./bin && cd bin && zip lambdaMailer.zip lambdaMailer templates
