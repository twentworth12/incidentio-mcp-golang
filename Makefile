.PHONY: build run test clean deps test-api

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
	@echo "Testing Incident.io API connection..."
	@if [ -f .env ]; then \
		export $$(cat .env | xargs) && \
		curl -H "Authorization: Bearer $$INCIDENT_IO_API_KEY" \
		     -H "Content-Type: application/json" \
		     "https://api.incident.io/v2/incidents?page_size=1" | jq .; \
	else \
		echo "No .env file found. Please create one with INCIDENT_IO_API_KEY=your-key"; \
	fi