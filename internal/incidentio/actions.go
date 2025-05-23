package incidentio

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// ListActionsOptions represents options for listing actions
type ListActionsOptions struct {
	PageSize   int
	After      string
	IncidentID string
	Status     []string
}

// ListActionsResponse represents the response from listing actions
type ListActionsResponse struct {
	Actions []Action `json:"actions"`
	ListResponse
}

// ListActions retrieves a list of actions
func (c *Client) ListActions(opts *ListActionsOptions) (*ListActionsResponse, error) {
	params := url.Values{}
	
	if opts != nil {
		if opts.PageSize > 0 {
			params.Set("page_size", strconv.Itoa(opts.PageSize))
		}
		if opts.After != "" {
			params.Set("after", opts.After)
		}
		if opts.IncidentID != "" {
			params.Set("incident_id", opts.IncidentID)
		}
		for _, status := range opts.Status {
			params.Add("status", status)
		}
	}

	respBody, err := c.doRequest("GET", "/actions", params, nil)
	if err != nil {
		return nil, err
	}

	var response ListActionsResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// GetAction retrieves a specific action by ID
func (c *Client) GetAction(id string) (*Action, error) {
	respBody, err := c.doRequest("GET", fmt.Sprintf("/actions/%s", id), nil, nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Action Action `json:"action"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response.Action, nil
}