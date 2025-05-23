package tools

import (
	"encoding/json"
	"fmt"

	"github.com/tomwentworth/incidentio-mcp-golang/internal/incidentio"
)

// ListIncidentTypesTool lists available incident types
type ListIncidentTypesTool struct {
	client *incidentio.Client
}

func NewListIncidentTypesTool(client *incidentio.Client) *ListIncidentTypesTool {
	return &ListIncidentTypesTool{client: client}
}

func (t *ListIncidentTypesTool) Name() string {
	return "list_incident_types"
}

func (t *ListIncidentTypesTool) Description() string {
	return "List available incident types"
}

func (t *ListIncidentTypesTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type":       "object",
		"properties": map[string]interface{}{},
		"additionalProperties": false,
	}
}

func (t *ListIncidentTypesTool) Execute(args map[string]interface{}) (string, error) {
	result, err := t.client.ListIncidentTypes()
	if err != nil {
		return "", fmt.Errorf("failed to list incident types: %w", err)
	}

	// Format the output to be more readable
	output := fmt.Sprintf("Found %d incident types:\n\n", len(result.IncidentTypes))
	
	for _, incidentType := range result.IncidentTypes {
		output += fmt.Sprintf("ID: %s\n", incidentType.ID)
		output += fmt.Sprintf("Name: %s\n", incidentType.Name)
		if incidentType.Description != "" {
			output += fmt.Sprintf("Description: %s\n", incidentType.Description)
		}
		if incidentType.IsDefault {
			output += "Default: Yes\n"
		}
		output += "\n"
	}

	// Also return the raw JSON
	jsonOutput, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return output, nil
	}

	return output + "\nRaw JSON:\n" + string(jsonOutput), nil
}