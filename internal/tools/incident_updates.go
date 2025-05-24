package tools

import (
	"encoding/json"
	"fmt"

	"github.com/tomwentworth/incidentio-mcp-golang/internal/incidentio"
)

// ListIncidentUpdatesTool lists incident updates
type ListIncidentUpdatesTool struct {
	client *incidentio.Client
}

func NewListIncidentUpdatesTool(client *incidentio.Client) *ListIncidentUpdatesTool {
	return &ListIncidentUpdatesTool{client: client}
}

func (t *ListIncidentUpdatesTool) Name() string {
	return "list_incident_updates"
}

func (t *ListIncidentUpdatesTool) Description() string {
	return "List incident updates (status messages posted during an incident)"
}

func (t *ListIncidentUpdatesTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"incident_id": map[string]interface{}{
				"type":        "string",
				"description": "Filter updates by incident ID",
			},
			"page_size": map[string]interface{}{
				"type":        "integer",
				"description": "Number of results per page (max 250)",
				"default":     25,
			},
		},
		"additionalProperties": false,
	}
}

func (t *ListIncidentUpdatesTool) Execute(args map[string]interface{}) (string, error) {
	opts := &incidentio.ListIncidentUpdatesOptions{}
	
	if incidentID, ok := args["incident_id"].(string); ok {
		opts.IncidentID = incidentID
	}
	if pageSize, ok := args["page_size"].(float64); ok {
		opts.PageSize = int(pageSize)
	}
	
	resp, err := t.client.ListIncidentUpdates(opts)
	if err != nil {
		return "", err
	}
	
	result, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to format response: %w", err)
	}
	
	return string(result), nil
}

// GetIncidentUpdateTool gets a specific incident update
type GetIncidentUpdateTool struct {
	client *incidentio.Client
}

func NewGetIncidentUpdateTool(client *incidentio.Client) *GetIncidentUpdateTool {
	return &GetIncidentUpdateTool{client: client}
}

func (t *GetIncidentUpdateTool) Name() string {
	return "get_incident_update"
}

func (t *GetIncidentUpdateTool) Description() string {
	return "Get details of a specific incident update by ID"
}

func (t *GetIncidentUpdateTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "string",
				"description": "The incident update ID",
			},
		},
		"required":             []interface{}{"id"},
		"additionalProperties": false,
	}
}

func (t *GetIncidentUpdateTool) Execute(args map[string]interface{}) (string, error) {
	id, ok := args["id"].(string)
	if !ok || id == "" {
		return "", fmt.Errorf("id parameter is required")
	}
	
	update, err := t.client.GetIncidentUpdate(id)
	if err != nil {
		return "", err
	}
	
	result, err := json.MarshalIndent(update, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to format response: %w", err)
	}
	
	return string(result), nil
}

// CreateIncidentUpdateTool creates a new incident update
type CreateIncidentUpdateTool struct {
	client *incidentio.Client
}

func NewCreateIncidentUpdateTool(client *incidentio.Client) *CreateIncidentUpdateTool {
	return &CreateIncidentUpdateTool{client: client}
}

func (t *CreateIncidentUpdateTool) Name() string {
	return "create_incident_update"
}

func (t *CreateIncidentUpdateTool) Description() string {
	return "Create a new incident update (status message) for an incident"
}

func (t *CreateIncidentUpdateTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"incident_id": map[string]interface{}{
				"type":        "string",
				"description": "The incident ID to post the update to",
			},
			"message": map[string]interface{}{
				"type":        "string",
				"description": "The update message to post",
			},
		},
		"required":             []interface{}{"incident_id", "message"},
		"additionalProperties": false,
	}
}

func (t *CreateIncidentUpdateTool) Execute(args map[string]interface{}) (string, error) {
	incidentID, ok := args["incident_id"].(string)
	if !ok || incidentID == "" {
		return "", fmt.Errorf("incident_id parameter is required")
	}
	
	message, ok := args["message"].(string)
	if !ok || message == "" {
		return "", fmt.Errorf("message parameter is required")
	}
	
	req := &incidentio.CreateIncidentUpdateRequest{
		IncidentID: incidentID,
		Message:    message,
	}
	
	update, err := t.client.CreateIncidentUpdate(req)
	if err != nil {
		return "", err
	}
	
	result, err := json.MarshalIndent(update, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to format response: %w", err)
	}
	
	return string(result), nil
}

// DeleteIncidentUpdateTool deletes an incident update
type DeleteIncidentUpdateTool struct {
	client *incidentio.Client
}

func NewDeleteIncidentUpdateTool(client *incidentio.Client) *DeleteIncidentUpdateTool {
	return &DeleteIncidentUpdateTool{client: client}
}

func (t *DeleteIncidentUpdateTool) Name() string {
	return "delete_incident_update"
}

func (t *DeleteIncidentUpdateTool) Description() string {
	return "Delete an incident update"
}

func (t *DeleteIncidentUpdateTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "string",
				"description": "The incident update ID to delete",
			},
		},
		"required":             []interface{}{"id"},
		"additionalProperties": false,
	}
}

func (t *DeleteIncidentUpdateTool) Execute(args map[string]interface{}) (string, error) {
	id, ok := args["id"].(string)
	if !ok || id == "" {
		return "", fmt.Errorf("id parameter is required")
	}
	
	if err := t.client.DeleteIncidentUpdate(id); err != nil {
		return "", err
	}
	
	return fmt.Sprintf("Successfully deleted incident update %s", id), nil
}