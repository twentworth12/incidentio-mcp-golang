#!/bin/bash

# Source environment variables
export INCIDENT_IO_API_KEY=inc_4f78bf72eb22ce2f6a58be0feb397a4030e1511c8d0eec234759ef83c00e8690

echo "Testing Incident.io API connection..."

# Test API connection with curl
curl -H "Authorization: Bearer $INCIDENT_IO_API_KEY" \
     -H "Content-Type: application/json" \
     https://api.incident.io/v2/incidents?page_size=1

echo -e "\n\nBuilding MCP server..."

# Build the server
go mod download
go mod tidy
go build -o bin/mcp-server cmd/mcp-server/main.go

echo "Build complete! Server binary created at bin/mcp-server"
echo "To run: export INCIDENT_IO_API_KEY=$INCIDENT_IO_API_KEY && ./bin/mcp-server"