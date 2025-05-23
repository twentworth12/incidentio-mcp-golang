package tools

import (
	"encoding/json"
	"fmt"

	"github.com/tomwentworth/incidentio-mcp-golang/internal/incidentio"
)

// ListSeveritiesTool lists available severities
type ListSeveritiesTool struct {
	client *incidentio.Client
}

func NewListSeveritiesTool(client *incidentio.Client) *ListSeveritiesTool {
	return &ListSeveritiesTool{client: client}
}

func (t *ListSeveritiesTool) Name() string {
	return "list_severities"
}

func (t *ListSeveritiesTool) Description() string {
	return "List available severity levels"
}

func (t *ListSeveritiesTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type":       "object",
		"properties": map[string]interface{}{},
		"additionalProperties": false,
	}
}

func (t *ListSeveritiesTool) Execute(args map[string]interface{}) (string, error) {
	result, err := t.client.ListSeverities()
	if err != nil {
		return "", fmt.Errorf("failed to list severities: %w", err)
	}

	// Format the output to be more readable
	output := fmt.Sprintf("Found %d severity levels:\n\n", len(result.Severities))
	
	for _, severity := range result.Severities {
		output += fmt.Sprintf("ID: %s\n", severity.ID)
		output += fmt.Sprintf("Name: %s\n", severity.Name)
		if severity.Description != "" {
			output += fmt.Sprintf("Description: %s\n", severity.Description)
		}
		output += fmt.Sprintf("Rank: %d\n", severity.Rank)
		output += "\n"
	}

	// Also return the raw JSON
	jsonOutput, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return output, nil
	}

	return output + "\nRaw JSON:\n" + string(jsonOutput), nil
}