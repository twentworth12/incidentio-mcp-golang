#!/bin/bash
set -e
cd /Users/tomwentworth/incidentio-mcp-golang
export INCIDENT_IO_API_KEY=inc_4f78bf72eb22ce2f6a58be0feb397a4030e1511c8d0eec234759ef83c00e8690

# Ensure the binary exists
if [ ! -f "./bin/mcp-server-clean" ]; then
    echo "Building MCP server..." >&2
    go build -o bin/mcp-server-clean cmd/mcp-server-clean/main.go
fi

# Run the server
exec ./bin/mcp-server-clean