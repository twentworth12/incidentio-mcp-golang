#!/bin/bash
# Wrapper script to capture debug logs
LOG_FILE="/tmp/mcp-incidentio-debug-$(date +%Y%m%d).log"
echo "$(date): MCP Server Starting" >> "$LOG_FILE"
exec /Users/tomwentworth/incidentio-mcp-golang/bin/mcp-server-clean 2>>"$LOG_FILE"