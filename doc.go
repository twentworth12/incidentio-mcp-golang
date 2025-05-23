/*
Package incidentio-mcp-golang provides a Model Context Protocol (MCP) server
implementation for interacting with the incident.io API.

This package implements a complete MCP server that exposes incident.io's
functionality through standardized tools that can be used by AI assistants
and other MCP clients.

Key Features:
  - Full incident.io V2 API coverage
  - Incident management (create, update, list, close)
  - Workflow automation and management
  - Alert routing and event handling
  - Role and user management
  - Comprehensive error handling

Usage:

	export INCIDENT_IO_API_KEY=your-api-key
	./bin/mcp-server

The server communicates via JSON-RPC over stdin/stdout, following the
MCP specification for tool discovery and execution.

For more information, see: https://github.com/twentworth12/incidentio-mcp-golang
*/
package main