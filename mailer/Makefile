include .envrc
export

# generate grpc code
.PHONY: compile/proto
compile/proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pb/mailer.proto

# build the lambda mailer (make sure to export .envrc)
.PHONY: build/lambda
build/lambda:
	GOOS=linux GOARCH=amd64 go build -o bin/lambdaMailer ./cmd/lambda && cp -r templates ./bin && cd bin && zip lambdaMailer.zip lambdaMailer templates
