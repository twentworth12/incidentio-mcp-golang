package incidentio

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// ListIncidentUpdates retrieves incident updates with optional filtering
func (c *Client) ListIncidentUpdates(opts *ListIncidentUpdatesOptions) (*ListIncidentUpdatesResponse, error) {
	params := url.Values{}
	
	if opts != nil {
		if opts.IncidentID != "" {
			params.Set("incident_id", opts.IncidentID)
		}
		if opts.PageSize > 0 {
			params.Set("page_size", strconv.Itoa(opts.PageSize))
		}
		if opts.After != "" {
			params.Set("after", opts.After)
		}
	}
	
	respBody, err := c.doRequest("GET", "/incident_updates", params, nil)
	if err != nil {
		return nil, err
	}
	
	var response ListIncidentUpdatesResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	return &response, nil
}

// GetIncidentUpdate retrieves a specific incident update by ID
func (c *Client) GetIncidentUpdate(id string) (*IncidentUpdate, error) {
	respBody, err := c.doRequest("GET", fmt.Sprintf("/incident_updates/%s", id), nil, nil)
	if err != nil {
		return nil, err
	}
	
	var response struct {
		IncidentUpdate IncidentUpdate `json:"incident_update"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	return &response.IncidentUpdate, nil
}

// CreateIncidentUpdate creates a new incident update
func (c *Client) CreateIncidentUpdate(req *CreateIncidentUpdateRequest) (*IncidentUpdate, error) {
	// Validate required fields
	if req.IncidentID == "" {
		return nil, fmt.Errorf("incident_id is required")
	}
	if req.Message == "" {
		return nil, fmt.Errorf("message is required")
	}
	
	respBody, err := c.doRequest("POST", "/incident_updates", nil, req)
	if err != nil {
		return nil, err
	}
	
	var response struct {
		IncidentUpdate IncidentUpdate `json:"incident_update"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	return &response.IncidentUpdate, nil
}

// DeleteIncidentUpdate deletes an incident update
func (c *Client) DeleteIncidentUpdate(id string) error {
	_, err := c.doRequest("DELETE", fmt.Sprintf("/incident_updates/%s", id), nil, nil)
	return err
}