#!/bin/bash
# Simple debug script that logs to a writable location
LOG_FILE="$HOME/mcp-debug.log"
echo "$(date): MCP Server Starting" >> "$LOG_FILE"
echo "PWD: $(pwd)" >> "$LOG_FILE"
echo "Script path: $0" >> "$LOG_FILE"
exec /Users/tomwentworth/incidentio-mcp-golang/bin/mcp-debug 2>>"$LOG_FILE"