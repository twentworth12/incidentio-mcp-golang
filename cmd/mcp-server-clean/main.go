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
	"github.com/tomwentworth/incidentio-mcp-golang/pkg/mcp"
)

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
		tools: make(map[string]tools.Tool),
	}
	server.registerTools()
	server.start(ctx)
}

type MCPServer struct {
	tools map[string]tools.Tool
}

func (s *MCPServer) registerTools() {
	// Try to initialize incident.io client
	client, err := incidentio.NewClient()
	if err != nil {
		// If client initialization fails, register only example tool
		exampleTool := &tools.ExampleTool{}
		s.tools[exampleTool.Name()] = exampleTool
		return
	}

	// Register all incident.io tools
	s.tools["list_incidents"] = tools.NewListIncidentsTool(client)
	s.tools["get_incident"] = tools.NewGetIncidentTool(client)
	s.tools["create_incident"] = tools.NewCreateIncidentTool(client)
	s.tools["update_incident"] = tools.NewUpdateIncidentTool(client)
	s.tools["close_incident"] = tools.NewCloseIncidentTool(client)
	s.tools["list_incident_statuses"] = tools.NewListIncidentStatusesTool(client)
	s.tools["list_alerts"] = tools.NewListAlertsTool(client)
	s.tools["get_alert"] = tools.NewGetAlertTool(client)
	s.tools["list_alerts_for_incident"] = tools.NewListAlertsForIncidentTool(client)
	s.tools["list_actions"] = tools.NewListActionsTool(client)
	s.tools["get_action"] = tools.NewGetActionTool(client)
	s.tools["list_available_incident_roles"] = tools.NewListIncidentRolesTool(client)
	s.tools["list_users"] = tools.NewListUsersTool(client)
	s.tools["assign_incident_role"] = tools.NewAssignIncidentRoleTool(client)
}

func (s *MCPServer) start(ctx context.Context) {
	encoder := json.NewEncoder(os.Stdout)
	decoder := json.NewDecoder(os.Stdin)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			var rawMsg json.RawMessage
			if err := decoder.Decode(&rawMsg); err != nil {
				if err == io.EOF {
					return
				}
				// Skip malformed JSON silently
				continue
			}

			// Try to parse as a proper JSON-RPC message
			var msg mcp.Message
			if err := json.Unmarshal(rawMsg, &msg); err != nil {
				// If we can't parse it, try to extract an ID to send proper error
				var partialMsg struct {
					ID      interface{} `json:"id"`
					Jsonrpc string      `json:"jsonrpc"`
				}
				if json.Unmarshal(rawMsg, &partialMsg) == nil && partialMsg.Jsonrpc == "2.0" {
					errorResp := &mcp.Message{
						Jsonrpc: "2.0",
						ID:      partialMsg.ID,
						Error: &mcp.Error{
							Code:    -32700,
							Message: "Parse error",
						},
					}
					encoder.Encode(errorResp)
				}
				continue
			}

			// Validate required fields
			if msg.Jsonrpc != "2.0" {
				if msg.ID != nil {
					errorResp := &mcp.Message{
						Jsonrpc: "2.0",
						ID:      msg.ID,
						Error: &mcp.Error{
							Code:    -32600,
							Message: "Invalid Request: missing or invalid jsonrpc field",
						},
					}
					encoder.Encode(errorResp)
				}
				continue
			}

			// Handle notifications (no ID) without response
			if msg.ID == nil {
				if msg.Method == "initialized" || msg.Method == "$/cancelled" {
					continue
				}
				// Unknown notification, ignore silently
				continue
			}

			response := s.handleMessage(&msg)
			if response != nil {
				encoder.Encode(response)
			}
		}
	}
}

func (s *MCPServer) handleMessage(msg *mcp.Message) *mcp.Message {
	// Ensure we always have an ID for responses (except notifications)
	if msg.ID == nil {
		return nil // This is a notification, no response needed
	}

	switch msg.Method {
	case "initialize":
		return &mcp.Message{
			Jsonrpc: "2.0",
			ID:      msg.ID,
			Result: map[string]interface{}{
				"protocolVersion": "2024-11-05",
				"capabilities": map[string]interface{}{
					"tools": map[string]interface{}{
						"listChanged": false,
					},
				},
				"serverInfo": map[string]interface{}{
					"name":    "incidentio-mcp-server",
					"version": "1.0.0",
				},
			},
		}
	case "initialized":
		// This should be handled as notification (no ID), but just in case
		return nil
	case "tools/list":
		var toolsList []map[string]interface{}
		for _, tool := range s.tools {
			toolsList = append(toolsList, map[string]interface{}{
				"name":        tool.Name(),
				"description": tool.Description(),
				"inputSchema": tool.InputSchema(),
			})
		}
		return &mcp.Message{
			Jsonrpc: "2.0",
			ID:      msg.ID,
			Result: map[string]interface{}{
				"tools": toolsList,
			},
		}
	case "tools/call":
		return s.handleToolCall(msg)
	default:
		return &mcp.Message{
			Jsonrpc: "2.0",
			ID:      msg.ID,
			Error: &mcp.Error{
				Code:    -32601,
				Message: fmt.Sprintf("Method not found: %s", msg.Method),
			},
		}
	}
}

func (s *MCPServer) handleToolCall(msg *mcp.Message) *mcp.Message {
	params, ok := msg.Params.(map[string]interface{})
	if !ok {
		return &mcp.Message{
			Jsonrpc: "2.0",
			ID:      msg.ID,
			Error: &mcp.Error{
				Code:    -32602,
				Message: "Invalid params",
			},
		}
	}

	toolName, ok := params["name"].(string)
	if !ok {
		return &mcp.Message{
			Jsonrpc: "2.0",
			ID:      msg.ID,
			Error: &mcp.Error{
				Code:    -32602,
				Message: "Missing tool name",
			},
		}
	}

	tool, exists := s.tools[toolName]
	if !exists {
		return &mcp.Message{
			Jsonrpc: "2.0",
			ID:      msg.ID,
			Error: &mcp.Error{
				Code:    -32602,
				Message: fmt.Sprintf("Tool not found: %s", toolName),
			},
		}
	}

	args, _ := params["arguments"].(map[string]interface{})
	result, err := tool.Execute(args)
	if err != nil {
		return &mcp.Message{
			Jsonrpc: "2.0",
			ID:      msg.ID,
			Error: &mcp.Error{
				Code:    -32603,
				Message: err.Error(),
			},
		}
	}

	return &mcp.Message{
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