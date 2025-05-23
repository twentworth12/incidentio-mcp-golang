package tools

import (
	"encoding/json"
	"fmt"

	"github.com/tomwentworth/incidentio-mcp-golang/internal/incidentio"
)

// ListActionsTool lists actions from Incident.io
type ListActionsTool struct {
	client *incidentio.Client
}

func NewListActionsTool(client *incidentio.Client) *ListActionsTool {
	return &ListActionsTool{client: client}
}

func (t *ListActionsTool) Name() string {
	return "list_actions"
}

func (t *ListActionsTool) Description() string {
	return "List actions from Incident.io with optional filters"
}

func (t *ListActionsTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"page_size": map[string]interface{}{
				"type":        "integer",
				"description": "Number of results per page (max 250)",
				"default":     25,
			},
			"incident_id": map[string]interface{}{
				"type":        "string",
				"description": "Filter actions by incident ID",
			},
			"status": map[string]interface{}{
				"type":        "array",
				"items":       map[string]interface{}{"type": "string"},
				"description": "Filter by action status (outstanding, completed, deleted)",
			},
		},
	}
}

func (t *ListActionsTool) Execute(args map[string]interface{}) (string, error) {
	opts := &incidentio.ListActionsOptions{}
	
	if pageSize, ok := args["page_size"].(float64); ok {
		opts.PageSize = int(pageSize)
	}
	
	if incidentID, ok := args["incident_id"].(string); ok {
		opts.IncidentID = incidentID
	}
	
	if statuses, ok := args["status"].([]interface{}); ok {
		for _, s := range statuses {
			if str, ok := s.(string); ok {
				opts.Status = append(opts.Status, str)
			}
		}
	}

	resp, err := t.client.ListActions(opts)
	if err != nil {
		return "", err
	}

	result, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to format response: %w", err)
	}

	return string(result), nil
}

// GetActionTool retrieves a specific action
type GetActionTool struct {
	client *incidentio.Client
}

func NewGetActionTool(client *incidentio.Client) *GetActionTool {
	return &GetActionTool{client: client}
}

func (t *GetActionTool) Name() string {
	return "get_action"
}

func (t *GetActionTool) Description() string {
	return "Get details of a specific action by ID"
}

func (t *GetActionTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "string",
				"description": "The action ID",
			},
		},
		"required": []string{"id"},
	}
}

func (t *GetActionTool) Execute(args map[string]interface{}) (string, error) {
	id, ok := args["id"].(string)
	if !ok {
		return "", fmt.Errorf("id parameter is required")
	}

	action, err := t.client.GetAction(id)
	if err != nil {
		return "", err
	}

	result, err := json.MarshalIndent(action, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to format response: %w", err)
	}

	return string(result), nil
}