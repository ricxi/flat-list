include .envrc
export

.PHONY: run/dev/token/http
run/dev/token/http:
	@echo "starting http server"
	go run ./cmd/http

.PHONY: run/dev/token/grpc
run/dev/token/grpc:
	@echo "starting grpc server"
	go run ./cmd/grpc
