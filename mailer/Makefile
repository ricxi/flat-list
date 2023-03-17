include .env
export

.PHONY: run/dev/http
run/dev/http:
	go run ./cmd/http


.PHONY: run/dev/grpc
run/dev/grpc:
	go run ./cmd/grpc

.PHONY: compile/proto
compile/proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pb/mailer.proto


# build the lambda mailer
.PHONY: build/lambda
build/lambda:
	GOOS=linux GOARCH=amd64 go build -o bin/lambdaMailer ./cmd/lambda && cp -r templates ./bin && cd bin && zip lambdaMailer.zip lambdaMailer templates
