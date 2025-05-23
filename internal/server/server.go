package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/tomwentworth/incidentio-mcp-golang/internal/incidentio"
	"github.com/tomwentworth/incidentio-mcp-golang/internal/tools"
	"github.com/tomwentworth/incidentio-mcp-golang/pkg/mcp"
)

type Server struct {
	tools map[string]tools.Tool
}

func New() *Server {
	return &Server{
		tools: make(map[string]tools.Tool),
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.registerTools()

	encoder := json.NewEncoder(os.Stdout)
	decoder := json.NewDecoder(os.Stdin)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			var msg mcp.Message
			if err := decoder.Decode(&msg); err != nil {
				if err == io.EOF {
					return nil
				}
				fmt.Fprintf(os.Stderr, "DEBUG: Failed to decode message: %v\n", err)
				continue
			}

			// Debug: log the incoming message
			msgBytes, _ := json.Marshal(msg)
			fmt.Fprintf(os.Stderr, "DEBUG: Received message: %s\n", string(msgBytes))

			response, err := s.handleMessage(&msg)
			if err != nil {
				response = s.createErrorResponse(msg.ID, err)
			}

			if response != nil {
				// Debug: log the response being sent
				responseBytes, _ := json.Marshal(response)
				fmt.Fprintf(os.Stderr, "DEBUG: Sending response: %s\n", string(responseBytes))
				encoder.Encode(response)
			}
		}
	}
}

func (s *Server) registerTools() {
	// Always register example tool
	exampleTool := &tools.ExampleTool{}
	s.tools[exampleTool.Name()] = exampleTool

	// Initialize Incident.io client
	client, err := incidentio.NewClient()
	if err != nil {
		// If client initialization fails, we still have the example tool
		return
	}

	// Register Incident tools
	s.tools["list_incidents"] = tools.NewListIncidentsTool(client)
	s.tools["get_incident"] = tools.NewGetIncidentTool(client)
	s.tools["create_incident"] = tools.NewCreateIncidentTool(client)
	s.tools["update_incident"] = tools.NewUpdateIncidentTool(client)
	s.tools["close_incident"] = tools.NewCloseIncidentTool(client)
	s.tools["list_incident_statuses"] = tools.NewListIncidentStatusesTool(client)

	// Register Alert tools
	s.tools["list_alerts"] = tools.NewListAlertsTool(client)
	s.tools["get_alert"] = tools.NewGetAlertTool(client)
	s.tools["list_alerts_for_incident"] = tools.NewListAlertsForIncidentTool(client)

	// Register Action tools
	s.tools["list_actions"] = tools.NewListActionsTool(client)
	s.tools["get_action"] = tools.NewGetActionTool(client)

	// Register Role tools
	s.tools["list_available_incident_roles"] = tools.NewListIncidentRolesTool(client)
	s.tools["list_users"] = tools.NewListUsersTool(client)
	s.tools["assign_incident_role"] = tools.NewAssignIncidentRoleTool(client)
}

func (s *Server) handleMessage(msg *mcp.Message) (*mcp.Message, error) {
	// Handle notifications (no ID means it's a notification)
	if msg.ID == nil {
		// Notifications don't require a response
		return nil, nil
	}

	switch msg.Method {
	case "initialize":
		return s.handleInitialize(msg)
	case "tools/list":
		return s.handleToolsList(msg)
	case "tools/call":
		return s.handleToolCall(msg)
	default:
		// Return proper JSON-RPC error for unknown methods
		return &mcp.Message{
			Jsonrpc: "2.0",
			ID:      msg.ID,
			Error: &mcp.Error{
				Code:    -32601,
				Message: fmt.Sprintf("Method not found: %s", msg.Method),
			},
		}, nil
	}
}

func (s *Server) handleInitialize(msg *mcp.Message) (*mcp.Message, error) {
	response := &mcp.Message{
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
	return response, nil
}

func (s *Server) handleToolsList(msg *mcp.Message) (*mcp.Message, error) {
	var toolsList []map[string]interface{}
	for _, tool := range s.tools {
		toolsList = append(toolsList, map[string]interface{}{
			"name":        tool.Name(),
			"description": tool.Description(),
			"inputSchema": tool.InputSchema(),
		})
	}

	response := &mcp.Message{
		Jsonrpc: "2.0",
		ID:      msg.ID,
		Result: map[string]interface{}{
			"tools": toolsList,
		},
	}
	return response, nil
}

func (s *Server) handleToolCall(msg *mcp.Message) (*mcp.Message, error) {
	params, ok := msg.Params.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid params")
	}

	toolName, ok := params["name"].(string)
	if !ok {
		return nil, fmt.Errorf("missing tool name")
	}

	tool, exists := s.tools[toolName]
	if !exists {
		return nil, fmt.Errorf("tool not found: %s", toolName)
	}

	args, _ := params["arguments"].(map[string]interface{})
	result, err := tool.Execute(args)
	if err != nil {
		return nil, err
	}

	response := &mcp.Message{
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
	return response, nil
}

func (s *Server) createErrorResponse(id interface{}, err error) *mcp.Message {
	return &mcp.Message{
		Jsonrpc: "2.0",
		ID:      id,
		Error: &mcp.Error{
			Code:    -32603,
			Message: err.Error(),
		},
	}
}