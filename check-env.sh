#!/bin/bash
echo "Environment Check:" >&2
echo "USER: $USER" >&2
echo "HOME: $HOME" >&2
echo "PATH: $PATH" >&2
echo "PWD: $PWD" >&2
echo "ANTHROPIC_API_KEY: ${ANTHROPIC_API_KEY:+SET}" >&2
echo "INCIDENT_IO_API_KEY: ${INCIDENT_IO_API_KEY:+SET}" >&2
echo "Go version:" >&2
go version >&2
echo "Starting MCP server..." >&2
exec /Users/tomwentworth/incidentio-mcp-golang/bin/mcp-debug 2>>/tmp/mcp-debug.log