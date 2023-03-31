include .env
export

.PHONY: run/dev
run/dev:
	@echo "running task service in dev mode"
	go run ./cmd/http

.PHONY: test
test:
	@echo "testing app"
	go test -cover