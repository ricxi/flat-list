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

.PHONY: test/e2e
test/e2e:
	@echo "e2e test"
	go test -v ./cmd/http