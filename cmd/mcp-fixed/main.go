package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/tomwentworth/incidentio-mcp-golang/internal/incidentio"
	"github.com/tomwentworth/incidentio-mcp-golang/internal/tools"
)

type Message struct {
	Jsonrpc string      `json:"jsonrpc"`
	Method  string      `json:"method,omitempty"`
	Params  interface{} `json:"params,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   *ErrorResp  `json:"error,omitempty"`
	ID      interface{} `json:"id,omitempty"`
}

type ErrorResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Tool interface {
	Name() string
	Description() string
	InputSchema() map[string]interface{}
	Execute(args map[string]interface{}) (string, error)
}

type MCPServer struct {
	tools map[string]Tool
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		cancel()
	}()

	server := &MCPServer{
		tools: make(map[string]Tool),
	}
	server.registerTools()
	server.start(ctx)
}

func (s *MCPServer) registerTools() {
	// Always register the example tool first
	exampleTool := &tools.ExampleTool{}
	s.tools[exampleTool.Name()] = exampleTool

	// Try to initialize Incident.io client
	client, err := incidentio.NewClient()
	if err != nil {
		// If client initialization fails, only example tool is available
		return
	}

	// Register Incident.io tools with simplified schemas
	s.tools["list_incidents"] = &SimpleListIncidentsTool{client: client}
	s.tools["get_incident"] = &SimpleGetIncidentTool{client: client}
	s.tools["update_incident"] = &SimpleUpdateIncidentTool{client: client}
}

func (s *MCPServer) start(ctx context.Context) {
	encoder := json.NewEncoder(os.Stdout)
	decoder := json.NewDecoder(os.Stdin)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			var msg Message
			if err := decoder.Decode(&msg); err != nil {
				if err == io.EOF {
					return
				}
				continue
			}

			response := s.handleMessage(&msg)
			if response != nil {
				encoder.Encode(response)
			}
		}
	}
}

func (s *MCPServer) handleMessage(msg *Message) *Message {
	switch msg.Method {
	case "initialize":
		return &Message{
			Jsonrpc: "2.0",
			ID:      msg.ID,
			Result: map[string]interface{}{
				"protocolVersion": "2024-11-05",
				"capabilities": map[string]interface{}{
					"tools": map[string]interface{}{},
				},
				"serverInfo": map[string]interface{}{
					"name":    "incidentio-mcp-server",
					"version": "0.1.0",
				},
			},
		}
	case "initialized":
		return nil
	case "tools/list":
		var toolsList []map[string]interface{}
		for _, tool := range s.tools {
			schema := tool.InputSchema()
			// Ensure schema is always a valid object
			if schema == nil {
				schema = map[string]interface{}{
					"type":       "object",
					"properties": map[string]interface{}{},
				}
			}
			
			toolsList = append(toolsList, map[string]interface{}{
				"name":        tool.Name(),
				"description": tool.Description(),
				"inputSchema": schema,
			})
		}
		return &Message{
			Jsonrpc: "2.0",
			ID:      msg.ID,
			Result: map[string]interface{}{
				"tools": toolsList,
			},
		}
	case "tools/call":
		return s.handleToolCall(msg)
	default:
		return &Message{
			Jsonrpc: "2.0",
			ID:      msg.ID,
			Error: &ErrorResp{
				Code:    -32601,
				Message: fmt.Sprintf("Method not found: %s", msg.Method),
			},
		}
	}
}

func (s *MCPServer) handleToolCall(msg *Message) *Message {
	params, ok := msg.Params.(map[string]interface{})
	if !ok {
		return &Message{
			Jsonrpc: "2.0",
			ID:      msg.ID,
			Error: &ErrorResp{
				Code:    -32602,
				Message: "Invalid params",
			},
		}
	}

	toolName, ok := params["name"].(string)
	if !ok {
		return &Message{
			Jsonrpc: "2.0",
			ID:      msg.ID,
			Error: &ErrorResp{
				Code:    -32602,
				Message: "Missing tool name",
			},
		}
	}

	tool, exists := s.tools[toolName]
	if !exists {
		return &Message{
			Jsonrpc: "2.0",
			ID:      msg.ID,
			Error: &ErrorResp{
				Code:    -32602,
				Message: fmt.Sprintf("Tool not found: %s", toolName),
			},
		}
	}

	args, _ := params["arguments"].(map[string]interface{})
	if args == nil {
		args = make(map[string]interface{})
	}

	result, err := tool.Execute(args)
	if err != nil {
		return &Message{
			Jsonrpc: "2.0",
			ID:      msg.ID,
			Error: &ErrorResp{
				Code:    -32603,
				Message: err.Error(),
			},
		}
	}

	return &Message{
		Jsonrpc: "2.0",
		ID:      msg.ID,
		Result: map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": result,
				},
			},
		},
	}
}

// Simplified tools with strict schemas

type SimpleListIncidentsTool struct {
	client *incidentio.Client
}

func (t *SimpleListIncidentsTool) Name() string {
	return "list_incidents"
}

func (t *SimpleListIncidentsTool) Description() string {
	return "List incidents from Incident.io"
}

func (t *SimpleListIncidentsTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"page_size": map[string]interface{}{
				"type":        "integer",
				"description": "Number of results per page",
				"minimum":     1,
				"maximum":     250,
			},
		},
		"additionalProperties": false,
	}
}

func (t *SimpleListIncidentsTool) Execute(args map[string]interface{}) (string, error) {
	opts := &incidentio.ListIncidentsOptions{}
	
	if pageSize, ok := args["page_size"].(float64); ok {
		opts.PageSize = int(pageSize)
	}

	resp, err := t.client.ListIncidents(opts)
	if err != nil {
		return "", err
	}

	result, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to format response: %w", err)
	}

	return string(result), nil
}

type SimpleGetIncidentTool struct {
	client *incidentio.Client
}

func (t *SimpleGetIncidentTool) Name() string {
	return "get_incident"
}

func (t *SimpleGetIncidentTool) Description() string {
	return "Get details of a specific incident by ID"
}

func (t *SimpleGetIncidentTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "string",
				"description": "The incident ID",
				"minLength":   1,
			},
		},
		"required":             []string{"id"},
		"additionalProperties": false,
	}
}

func (t *SimpleGetIncidentTool) Execute(args map[string]interface{}) (string, error) {
	id, ok := args["id"].(string)
	if !ok || id == "" {
		return "", fmt.Errorf("id parameter is required and must be a non-empty string")
	}

	incident, err := t.client.GetIncident(id)
	if err != nil {
		return "", err
	}

	result, err := json.MarshalIndent(incident, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to format response: %w", err)
	}

	return string(result), nil
}

type SimpleUpdateIncidentTool struct {
	client *incidentio.Client
}

func (t *SimpleUpdateIncidentTool) Name() string {
	return "update_incident"
}

func (t *SimpleUpdateIncidentTool) Description() string {
	return "Update an existing incident status"
}

func (t *SimpleUpdateIncidentTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "string",
				"description": "The incident ID to update",
				"minLength":   1,
			},
			"incident_status_id": map[string]interface{}{
				"type":        "string",
				"description": "The new incident status ID",
				"minLength":   1,
			},
		},
		"required":             []string{"id", "incident_status_id"},
		"additionalProperties": false,
	}
}

func (t *SimpleUpdateIncidentTool) Execute(args map[string]interface{}) (string, error) {
	id, ok := args["id"].(string)
	if !ok || id == "" {
		return "", fmt.Errorf("id parameter is required and must be a non-empty string")
	}

	statusID, ok := args["incident_status_id"].(string)
	if !ok || statusID == "" {
		return "", fmt.Errorf("incident_status_id parameter is required and must be a non-empty string")
	}

	req := &incidentio.UpdateIncidentRequest{
		IncidentStatusID: statusID,
	}

	incident, err := t.client.UpdateIncident(id, req)
	if err != nil {
		return "", err
	}

	result, err := json.MarshalIndent(incident, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to format response: %w", err)
	}

	return string(result), nil
}