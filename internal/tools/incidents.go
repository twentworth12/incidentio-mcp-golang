package tools

import (
	"encoding/json"
	"fmt"

	"github.com/tomwentworth/incidentio-mcp-golang/internal/incidentio"
)

// ListIncidentsTool lists incidents from Incident.io
type ListIncidentsTool struct {
	client *incidentio.Client
}

func NewListIncidentsTool(client *incidentio.Client) *ListIncidentsTool {
	return &ListIncidentsTool{client: client}
}

func (t *ListIncidentsTool) Name() string {
	return "list_incidents"
}

func (t *ListIncidentsTool) Description() string {
	return "List incidents from Incident.io with optional filters"
}

func (t *ListIncidentsTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"page_size": map[string]interface{}{
				"type":        "integer",
				"description": "Number of results per page (max 250)",
				"default":     25,
			},
			"status": map[string]interface{}{
				"type":        "array",
				"items":       map[string]interface{}{"type": "string"},
				"description": "Filter by incident status (e.g., triage, active, resolved, closed)",
			},
			"severity": map[string]interface{}{
				"type":        "array",
				"items":       map[string]interface{}{"type": "string"},
				"description": "Filter by severity",
			},
		},
	}
}

func (t *ListIncidentsTool) Execute(args map[string]interface{}) (string, error) {
	opts := &incidentio.ListIncidentsOptions{}
	
	if pageSize, ok := args["page_size"].(float64); ok {
		opts.PageSize = int(pageSize)
	}
	
	if statuses, ok := args["status"].([]interface{}); ok {
		for _, s := range statuses {
			if str, ok := s.(string); ok {
				opts.Status = append(opts.Status, str)
			}
		}
	}
	
	if severities, ok := args["severity"].([]interface{}); ok {
		for _, s := range severities {
			if str, ok := s.(string); ok {
				opts.Severity = append(opts.Severity, str)
			}
		}
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

// GetIncidentTool retrieves a specific incident
type GetIncidentTool struct {
	client *incidentio.Client
}

func NewGetIncidentTool(client *incidentio.Client) *GetIncidentTool {
	return &GetIncidentTool{client: client}
}

func (t *GetIncidentTool) Name() string {
	return "get_incident"
}

func (t *GetIncidentTool) Description() string {
	return "Get details of a specific incident by ID"
}

func (t *GetIncidentTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "string",
				"description": "The incident ID",
			},
		},
		"required": []string{"id"},
	}
}

func (t *GetIncidentTool) Execute(args map[string]interface{}) (string, error) {
	id, ok := args["id"].(string)
	if !ok {
		return "", fmt.Errorf("id parameter is required")
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

// CreateIncidentTool creates a new incident
type CreateIncidentTool struct {
	client *incidentio.Client
}

func NewCreateIncidentTool(client *incidentio.Client) *CreateIncidentTool {
	return &CreateIncidentTool{client: client}
}

func (t *CreateIncidentTool) Name() string {
	return "create_incident"
}

func (t *CreateIncidentTool) Description() string {
	return "Create a new incident in Incident.io"
}

func (t *CreateIncidentTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "The incident name/title",
			},
			"summary": map[string]interface{}{
				"type":        "string",
				"description": "A summary of the incident",
			},
			"status": map[string]interface{}{
				"type":        "string",
				"description": "Initial status (triage, active, resolved, closed)",
				"default":     "triage",
			},
			"severity_id": map[string]interface{}{
				"type":        "string",
				"description": "The severity ID",
			},
			"incident_type_id": map[string]interface{}{
				"type":        "string",
				"description": "The incident type ID",
			},
			"incident_status_id": map[string]interface{}{
				"type":        "string",
				"description": "The incident status ID",
			},
			"slack_channel_name_override": map[string]interface{}{
				"type":        "string",
				"description": "Override the auto-generated Slack channel name",
			},
		},
		"required": []string{"name"},
	}
}

func (t *CreateIncidentTool) Execute(args map[string]interface{}) (string, error) {
	name, ok := args["name"].(string)
	if !ok {
		return "", fmt.Errorf("name parameter is required")
	}

	req := &incidentio.CreateIncidentRequest{
		Name: name,
	}

	if summary, ok := args["summary"].(string); ok {
		req.Summary = summary
	}
	if statusID, ok := args["incident_status_id"].(string); ok {
		req.IncidentStatusID = statusID
	}
	if severityID, ok := args["severity_id"].(string); ok {
		req.SeverityID = severityID
	}
	if typeID, ok := args["incident_type_id"].(string); ok {
		req.IncidentTypeID = typeID
	}
	if slackOverride, ok := args["slack_channel_name_override"].(string); ok {
		req.SlackChannelNameOverride = slackOverride
	}

	incident, err := t.client.CreateIncident(req)
	if err != nil {
		return "", err
	}

	result, err := json.MarshalIndent(incident, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to format response: %w", err)
	}

	return string(result), nil
}

// UpdateIncidentTool updates an existing incident
type UpdateIncidentTool struct {
	client *incidentio.Client
}

func NewUpdateIncidentTool(client *incidentio.Client) *UpdateIncidentTool {
	return &UpdateIncidentTool{client: client}
}

func (t *UpdateIncidentTool) Name() string {
	return "update_incident"
}

func (t *UpdateIncidentTool) Description() string {
	return "Update an existing incident"
}

func (t *UpdateIncidentTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "string",
				"description": "The incident ID to update",
			},
			"name": map[string]interface{}{
				"type":        "string",
				"description": "Update the incident name",
			},
			"summary": map[string]interface{}{
				"type":        "string",
				"description": "Update the incident summary",
			},
			"incident_status_id": map[string]interface{}{
				"type":        "string",
				"description": "Update the incident status ID",
			},
			"severity_id": map[string]interface{}{
				"type":        "string",
				"description": "Update the severity ID",
			},
		},
		"required": []string{"id"},
	}
}

func (t *UpdateIncidentTool) Execute(args map[string]interface{}) (string, error) {
	id, ok := args["id"].(string)
	if !ok {
		return "", fmt.Errorf("id parameter is required")
	}

	req := &incidentio.UpdateIncidentRequest{}
	hasUpdate := false

	if name, ok := args["name"].(string); ok {
		req.Name = name
		hasUpdate = true
	}
	if summary, ok := args["summary"].(string); ok {
		req.Summary = summary
		hasUpdate = true
	}
	if statusID, ok := args["incident_status_id"].(string); ok {
		req.IncidentStatusID = statusID
		hasUpdate = true
	}
	if severityID, ok := args["severity_id"].(string); ok {
		req.SeverityID = severityID
		hasUpdate = true
	}

	if !hasUpdate {
		return "", fmt.Errorf("at least one field to update must be provided")
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