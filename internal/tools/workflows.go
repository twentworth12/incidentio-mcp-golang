package tools

import (
	"encoding/json"
	"fmt"

	"github.com/tomwentworth/incidentio-mcp-golang/internal/incidentio"
)

// ListWorkflowsTool lists workflows from incident.io
type ListWorkflowsTool struct {
	client *incidentio.Client
}

func NewListWorkflowsTool(client *incidentio.Client) *ListWorkflowsTool {
	return &ListWorkflowsTool{client: client}
}

func (t *ListWorkflowsTool) Name() string {
	return "list_workflows"
}

func (t *ListWorkflowsTool) Description() string {
	return "List workflows from incident.io with optional pagination"
}

func (t *ListWorkflowsTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"page_size": map[string]interface{}{
				"type":        "integer",
				"description": "Number of results per page",
				"minimum":     1,
				"maximum":     250,
			},
			"after": map[string]interface{}{
				"type":        "string",
				"description": "Pagination cursor for next page",
			},
		},
		"additionalProperties": false,
	}
}

func (t *ListWorkflowsTool) Execute(args map[string]interface{}) (string, error) {
	params := &incidentio.ListWorkflowsParams{}
	
	if pageSize, ok := args["page_size"].(float64); ok {
		params.PageSize = int(pageSize)
	}
	if after, ok := args["after"].(string); ok {
		params.After = after
	}
	
	result, err := t.client.ListWorkflows(params)
	if err != nil {
		return "", fmt.Errorf("failed to list workflows: %w", err)
	}
	
	output, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}
	
	return string(output), nil
}

// GetWorkflowTool gets details of a specific workflow
type GetWorkflowTool struct {
	client *incidentio.Client
}

func NewGetWorkflowTool(client *incidentio.Client) *GetWorkflowTool {
	return &GetWorkflowTool{client: client}
}

func (t *GetWorkflowTool) Name() string {
	return "get_workflow"
}

func (t *GetWorkflowTool) Description() string {
	return "Get details of a specific workflow by ID"
}

func (t *GetWorkflowTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "string",
				"description": "The workflow ID",
				"minLength":   1,
			},
		},
		"required":             []string{"id"},
		"additionalProperties": false,
	}
}

func (t *GetWorkflowTool) Execute(args map[string]interface{}) (string, error) {
	id, ok := args["id"].(string)
	if !ok || id == "" {
		return "", fmt.Errorf("workflow ID is required")
	}
	
	workflow, err := t.client.GetWorkflow(id)
	if err != nil {
		return "", fmt.Errorf("failed to get workflow: %w", err)
	}
	
	output, err := json.MarshalIndent(workflow, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}
	
	return string(output), nil
}

// UpdateWorkflowTool updates a workflow
type UpdateWorkflowTool struct {
	client *incidentio.Client
}

func NewUpdateWorkflowTool(client *incidentio.Client) *UpdateWorkflowTool {
	return &UpdateWorkflowTool{client: client}
}

func (t *UpdateWorkflowTool) Name() string {
	return "update_workflow"
}

func (t *UpdateWorkflowTool) Description() string {
	return "Update a workflow's configuration"
}

func (t *UpdateWorkflowTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "string",
				"description": "The workflow ID to update",
				"minLength":   1,
			},
			"name": map[string]interface{}{
				"type":        "string",
				"description": "New name for the workflow",
			},
			"enabled": map[string]interface{}{
				"type":        "boolean",
				"description": "Whether the workflow should be enabled",
			},
			"state": map[string]interface{}{
				"type":        "object",
				"description": "State configuration for the workflow",
			},
		},
		"required":             []string{"id"},
		"additionalProperties": false,
	}
}

func (t *UpdateWorkflowTool) Execute(args map[string]interface{}) (string, error) {
	id, ok := args["id"].(string)
	if !ok || id == "" {
		return "", fmt.Errorf("workflow ID is required")
	}
	
	req := &incidentio.UpdateWorkflowRequest{}
	
	if name, ok := args["name"].(string); ok {
		req.Name = name
	}
	
	if enabled, ok := args["enabled"].(bool); ok {
		req.Enabled = &enabled
	}
	
	if state, ok := args["state"].(map[string]interface{}); ok {
		req.State = state
	}
	
	workflow, err := t.client.UpdateWorkflow(id, req)
	if err != nil {
		return "", fmt.Errorf("failed to update workflow: %w", err)
	}
	
	output, err := json.MarshalIndent(workflow, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}
	
	return string(output), nil
}