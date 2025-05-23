package incidentio

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// ListIncidentsOptions represents options for listing incidents
type ListIncidentsOptions struct {
	PageSize int
	After    string
	Status   []string
	Severity []string
}

// ListIncidentsResponse represents the response from listing incidents
type ListIncidentsResponse struct {
	Incidents []Incident `json:"incidents"`
	ListResponse
}

// ListIncidents retrieves a list of incidents
func (c *Client) ListIncidents(opts *ListIncidentsOptions) (*ListIncidentsResponse, error) {
	params := url.Values{}
	
	if opts != nil {
		if opts.PageSize > 0 {
			params.Set("page_size", strconv.Itoa(opts.PageSize))
		}
		if opts.After != "" {
			params.Set("after", opts.After)
		}
		for _, status := range opts.Status {
			params.Add("status", status)
		}
		for _, severity := range opts.Severity {
			params.Add("severity", severity)
		}
	}

	respBody, err := c.doRequest("GET", "/incidents", params, nil)
	if err != nil {
		return nil, err
	}

	var response ListIncidentsResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// GetIncident retrieves a specific incident by ID
func (c *Client) GetIncident(id string) (*Incident, error) {
	respBody, err := c.doRequest("GET", fmt.Sprintf("/incidents/%s", id), nil, nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Incident Incident `json:"incident"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response.Incident, nil
}

// CreateIncident creates a new incident
func (c *Client) CreateIncident(req *CreateIncidentRequest) (*Incident, error) {
	respBody, err := c.doRequest("POST", "/incidents", nil, req)
	if err != nil {
		return nil, err
	}

	var response struct {
		Incident Incident `json:"incident"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response.Incident, nil
}

// UpdateIncident updates an existing incident using V2 actions/edit API
func (c *Client) UpdateIncident(id string, req *UpdateIncidentRequest) (*Incident, error) {
	// Use the correct V2 actions/edit endpoint
	editRequest := map[string]interface{}{
		"notify_incident_channel": true,
	}
	
	// Build the incident object with only the fields that are being updated
	incident := make(map[string]interface{})
	
	if req.Name != "" {
		incident["name"] = req.Name
	}
	if req.Summary != "" {
		incident["summary"] = req.Summary
	}
	if req.IncidentStatusID != "" {
		incident["incident_status_id"] = req.IncidentStatusID
	}
	if req.SeverityID != "" {
		incident["severity_id"] = req.SeverityID
	}
	if len(req.IncidentRoleAssignments) > 0 {
		incident["incident_role_assignments"] = req.IncidentRoleAssignments
	}
	
	// Only include incident object if there are fields to update
	if len(incident) > 0 {
		editRequest["incident"] = incident
	} else {
		return nil, fmt.Errorf("no fields to update")
	}
	
	respBody, err := c.doRequest("POST", fmt.Sprintf("/incidents/%s/actions/edit", id), nil, editRequest)
	if err != nil {
		return nil, err
	}

	var response struct {
		Incident Incident `json:"incident"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response.Incident, nil
}