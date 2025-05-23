# incident.io MCP Server

A GoLang implementation of an MCP (Model Context Protocol) server for incident.io, providing tools to interact with the incident.io V2 API.

## Project Structure

```
.
├── cmd/mcp-server/      # Main application entry point
├── internal/            # Private application code
│   ├── server/          # MCP server implementation
│   └── tools/           # Tool implementations
├── pkg/mcp/             # MCP protocol types and utilities
├── go.mod               # Go module definition
├── go.sum               # Go module checksums
└── Makefile             # Build commands
```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- incident.io API key (set as `INCIDENT_IO_API_KEY` environment variable)

### Installation

1. Clone the repository
2. Install dependencies:
   ```bash
   make deps
   ```

### Building

```bash
make build
```

This will create a binary in the `bin/` directory.

### Running

Set your incident.io API key:
```bash
export INCIDENT_IO_API_KEY=your-api-key
```

Then run the server:
```bash
make run
```

Or after building:
```bash
./bin/mcp-server
```

### Testing

```bash
make test
```

## Adding New Tools

1. Create a new file in `internal/tools/` implementing the `Tool` interface
2. Register the tool in `server.registerTools()` method in `internal/server/server.go`

Example tool implementation:
```go
type MyTool struct{}

func (t *MyTool) Name() string {
    return "my_tool"
}

func (t *MyTool) Description() string {
    return "Description of what the tool does"
}

func (t *MyTool) InputSchema() map[string]interface{} {
    return map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            // Define your parameters here
        },
        "required": []string{/* required parameters */},
    }
}

func (t *MyTool) Execute(args map[string]interface{}) (string, error) {
    // Tool implementation
    return "result", nil
}
```

## Available Tools

### Incident Management
- `list_incidents` - List incidents with optional filters (status, severity)
- `get_incident` - Get details of a specific incident by ID
- `create_incident` - Create a new incident
- `update_incident` - Update an existing incident
- `close_incident` - Close an incident with proper workflow handling
- `list_incident_statuses` - List available incident statuses

### Alert Management
- `list_alerts` - List alerts with optional filters
- `get_alert` - Get details of a specific alert by ID
- `list_alerts_for_incident` - List alerts associated with a specific incident
- `list_alert_sources` - List available alert sources
- `create_alert_event` - Create an alert event

### Alert Routing
- `list_alert_routes` - List alert routes with optional pagination
- `get_alert_route` - Get details of a specific alert route
- `create_alert_route` - Create a new alert route with conditions and escalations
- `update_alert_route` - Update an alert route's configuration

### Workflow Management
- `list_workflows` - List workflows with optional pagination
- `get_workflow` - Get details of a specific workflow
- `update_workflow` - Update a workflow's configuration

### Action Management
- `list_actions` - List actions with optional filters (incident_id, status)
- `get_action` - Get details of a specific action by ID

### Roles and Users
- `list_available_incident_roles` - List available incident roles
- `list_users` - List users in the organization
- `assign_incident_role` - Assign a role to a user for an incident

### Testing
- `example_tool` - A simple echo tool for testing

## MCP Protocol

This server implements the Model Context Protocol (MCP) for communication with AI assistants. The server:
- Communicates via JSON-RPC over stdin/stdout
- Supports tool registration and execution
- Follows the MCP 2024-11-05 protocol version
- Integrates with incident.io V2 API endpoints

## Environment Variables

- `INCIDENT_IO_API_KEY` (required) - Your incident.io API key
- `INCIDENT_IO_BASE_URL` (optional) - Override the API base URL (defaults to https://api.incident.io/v2)