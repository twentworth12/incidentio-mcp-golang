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
	// Try to initialize Incident.io client
	client, err := incidentio.NewClient()
	if err != nil {
		// If client initialization fails, register only example tool
		exampleTool := &tools.ExampleTool{}
		s.tools[exampleTool.Name()] = exampleTool
		return
	}

	// Register all Incident.io tools
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
			var msg mcp.Message
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

func (s *MCPServer) handleMessage(msg *mcp.Message) *mcp.Message {
	switch msg.Method {
	case "initialize":
		return &mcp.Message{
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
		// Return nil for initialized - no response needed
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