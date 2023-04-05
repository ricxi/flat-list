
.PHONY: test/response
test/response:
	@echo "Testing the response package..."
	go test ./response
	@echo "Testing complete for response package."

.PHONY: test/request
test/request:
	@echo "Testing the request package..."
	go test ./request
	@echo "Testing complete for request package."

.PHONY: benchmark/response
benchmark/response:
	@echo "Benchmarking the request package..."
	go test -bench=. -benchmem -count=5 -benchtime=5s ./response | tee response_benchmarks.txt
	@echo "Benchmarks complete for request package."

.PHONY: benchmark/request
benchmark/request:
	@echo "Benchmarking the request package..."
	go test -bench=. -benchmem -count=5 -benchtime=5s ./request
	@echo "Benchmarks complete for request package."

.PHONY: log/benchmark/request
log/benchmark/request:
	@echo "Benchmarking the request package..."
	go test -bench=. -benchmem -count=5 -benchtime=5s ./request | tee request_benchmarks.txt
	@echo "Benchmarks complete for request package."