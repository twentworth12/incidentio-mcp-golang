.PHONY: build run test clean deps test-api test-unit test-integration

build:
	go build -o bin/mcp-server cmd/mcp-server/main.go

run:
	@if [ -f .env ]; then export $$(cat .env | xargs); fi && go run cmd/mcp-server/main.go

test:
	go test ./...

clean:
	rm -rf bin/

deps:
	go mod download
	go mod tidy

test-api:
	@echo "Testing incident.io API endpoints..."
	@./test_endpoints.sh

test-unit:
	@echo "Running unit tests..."
	go test -v ./internal/incidentio/...

test-integration: test-unit build test-api
	@echo "All integration tests completed!"