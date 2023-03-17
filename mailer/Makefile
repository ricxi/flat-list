include .env
export

.PHONY: run/dev/http
run/dev/http:
	go run ./cmd/http


.PHONY: run/dev/grpc
run/dev/grpc:
	go run ./cmd/grpc

# build the lambda mailer
.PHONY: build/lambda
build/lambda:
	GOOS=linux GOARCH=amd64 go build -o bin/lambdaMailer ./cmd/lambda && cp -r templates ./bin && cd bin && zip lambdaMailer.zip lambdaMailer templates